package scheduled_tasks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/service"
	"mproxy/utils/taskRunManager"
)

// 把子账号的流量从redis同步到MySQL
func (m *manager) updateDynamicAccountFlowFromRedisToDB(ctx context.Context) {
	timneOut := 5 * time.Second
	ticker := time.NewTicker(timneOut)
	defer ticker.Stop()
	trm := taskRunManager.NewChanTask(5)
	for {
		ticker.Reset(timneOut)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]把子账号的流量从redis同步到MySQL任务上下文结束退出")
			return
		case <-ticker.C:

			_, results, err := common.
				GetRedisDB().
				BZMPop(
					ctx,
					5*time.Second, // 等待超时5秒
					"MIN",
					30, // 弹出30个
					constant.DynamicAccountIDFlowRedisQueueSortedSet,
				).Result()
			if err != nil {
				if err != redis.Nil {
					log.Error("[定时任务scheduled_tasks]把子账号的流量从redis同步到MySQL执行BZMPop命令时出错", zap.Any("error", err))
				}

				continue
			}

			//  处理弹出的元素信息
			for _, item := range results {
				num, _ := strconv.ParseInt(fmt.Sprint(item.Member), 10, 64)
				if num > 0 {
					trm.Run(func() {
						///获取分布式锁
						lockkey := fmt.Sprint(constant.DynamicAccountDbLock, num)
						if !common.
							GetRedisDB().
							SetNX(
								context.Background(),
								lockkey,
								"", constant.DynamicAccountDbLockTtl,
							).Val() {
							return
						}

						defer common.GetRedisDB().Del(context.Background(), lockkey)

						err := service.UpdateDynamicAccountFlowFromRedisToDB(
							context.Background(),
							common.GetMysqlDB(),
							common.GetRedisDB(),
							num,
						)
						if err != nil {
							log.Error("[定时任务scheduled_tasks]把子账号的流量从redis同步到MySQL更新子账号错误", zap.Any("accountId", num), zap.Any("error", err))
							return
						}

						///更新子账号缓存
						if err := service.UpdateDynamicAccountDataCachebyRedisPipe(
							context.Background(),
							common.GetMysqlDB(),
							common.GetRedisDB(),
							num,
						); err != nil {
							log.Error("[定时任务scheduled_tasks]更新子账号缓存错误", zap.Any("accountId", num), zap.Any("error", err))
						}
					})
				}
			}
			trm.Wait()

		}
	}
}
