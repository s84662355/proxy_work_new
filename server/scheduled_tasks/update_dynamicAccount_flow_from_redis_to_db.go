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
)

// /把子账号的流量从redis同步到MySQL
func (m *manager) updateDynamicAccountFlowFromRedisToDB(ctx context.Context) {
	timneOut := 5 * time.Second
	ticker := time.NewTicker(timneOut)
	defer ticker.Stop()
	for {
		ticker.Reset(timneOut)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks] 定时执行 updateDynamicAccountFlowFromRedisToDB 上下文Done() 退出")
			return
		case <-ticker.C:

			_, results, err := common.
				GetRedisDB().
				BZMPop(
					ctx,
					5*time.Second, // 等待超时5秒
					"MIN ",
					30, // 弹出30个
					constant.DynamicAccountIDFlowRedisQueueSortedSet,
				).Result()
			if err != nil {
				if err != redis.Nil {
					log.Error("[定时任务scheduled_tasks] updateDynamicAccountFlowFromRedisToDB 执行 BZMPop 命令时出错 ", zap.Any("error", err))
				}

				continue
			}

			//  处理弹出的元素信息
			for _, item := range results {
				num, _ := strconv.ParseInt(fmt.Sprint(item.Member), 10, 64)
				if num > 0 {
					err := service.UpdateDynamicAccountFlowFromRedisToDB(
						context.Background(),
						common.GetMysqlDB(),
						common.GetRedisDB(),
						num,
					)
					if err != nil {
						log.Error("[定时任务scheduled_tasks] updateDynamicAccountFlowFromRedisToDB 更新子账号", zap.Any("error", err))
					}
				}
			}

		}
	}
}
