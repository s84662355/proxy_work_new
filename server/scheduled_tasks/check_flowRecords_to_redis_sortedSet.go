package scheduled_tasks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/config"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/service"
	"mproxy/utils/taskRunManager"
)

// /检查流量表
func (m *manager) checkFlowRecordsToRedisSortedSet(ctx context.Context) {
	loopTime := 60 * time.Second
	ticker := time.NewTicker(loopTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks] 定时执行 checkFlowRecordsToRedisSortedSet 上下文Done() 退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks] 定时执行 checkFlowRecordsToRedisSortedSet 开始 ", zap.Time("Time", t))
			v, err := service.CheckFlowRecordsToRedisSortedSet(
				ctx,
				common.GetMysqlDB(),
				common.GetRedisDB(),
			)
			if err != nil {
				log.Error("[定时任务scheduled_tasks] 定时执行 checkFlowRecordsToRedisSortedSet   执行错误", zap.Time("Time", t), zap.Any("error", err))
			}
			log.Info("[定时任务scheduled_tasks] 定时执行 checkFlowRecordsToRedisSortedSet 结束  ", zap.String("Time", fmt.Sprintf("%+v---%+v", t, time.Now())), zap.Int64("count", v))
			ticker.Reset(loopTime)
		}
	}
}

func (m *manager) updateFlowRecordsFromRedisSortedSet(ctx context.Context) {
	loopTime := 5 * time.Second
	ticker := time.NewTicker(loopTime)
	defer ticker.Stop()
	trm := taskRunManager.NewChanTask(5)
	for {
		ticker.Reset(loopTime)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks] 定时执行 getFlowRecordsFromRedisSortedSet 上下文Done() 退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks] 定时执行 getFlowRecordsFromRedisSortedSet 开始 ", zap.Time("Time", t))
			_, results, err := common.GetRedisDB().
				BZMPop(
					ctx,
					5*time.Second, // 等待超时5秒
					"MIN ",
					30, // 弹出30个
					constant.FlowUserIdQueueSortedSet,
				).Result()
			if err != nil {
				if err != redis.Nil {
					log.Error("[定时任务scheduled_tasks] getFlowRecordsFromRedisSortedSet 执行 BZMPop 命令时出错 ", zap.Any("error", err))
				}
				continue
			}

			if len(results) > 0 {
				//  处理弹出的元素信息
				for _, v := range results {
					members := strings.Split(fmt.Sprint(v.Member), ",")
					if len(members) == 0 {
						continue
					}
					userId, err := strconv.ParseInt(members[0], 10, 64)
					if err != nil {
						log.Error("[定时任务scheduled_tasks] getFlowRecordsFromRedisSortedSet strconv.ParseInt 出错 ", zap.Any("error", err))
						continue
					}

					var recordId uint64 = 0
					if len(members) > 1 {
						recordId, _ = strconv.ParseUint(members[1], 10, 64)
					}

					if userId != 0 {
						trm.Run(func() {
							m.updateAndAutoBuyFlow(userId, recordId)
						})
					}
				}
				trm.Wait()
			}

		}
	}
}

func (m *manager) updateAndAutoBuyFlow(
	userId int64,
	recordId uint64,
) {
	lockKey := fmt.Sprint(constant.VsIPTransitDynamicUpdateFlowRedisLock, userId)
	if !common.GetRedisDB().
		SetNX(
			context.Background(),
			lockKey,
			"",
			constant.VsIPTransitDynamicUpdateFlowRedisLockTtl,
		).Val() {
		return
	}
	defer common.GetRedisDB().Del(context.Background(), lockKey)
	if _, _, err := service.UpdateFlowRecordsToDynamicFlow(
		context.Background(),
		common.GetMysqlDB(),
		userId,
		config.GetConf().FlowIncRate,
		300,
		recordId,
	); err != nil {
		log.Error("[定时任务scheduled_tasks] getFlowRecordsFromRedisSortedSet 更新流量出错 ", zap.Int64("userId", userId), zap.Any("error", err))
		return
	}

	if !common.GetRedisDB().
		SetNX(
			context.Background(),
			fmt.Sprint(constant.VsIPTransitDynamicAutoBuyFlowRedisLock, userId),
			"",
			constant.VsIPTransitDynamicAutoBuyFlowRedisLockTtl,
		).Val() {
		return
	}
	if err := service.DynamicFlowAndAutoBuy(
		context.Background(),
		common.GetMysqlDB(),
		userId,
	); err != nil {
		log.Error("[定时任务scheduled_tasks] getFlowRecordsFromRedisSortedSet 自动购买流量出错 ", zap.Int64("userId", userId), zap.Any("error", err))
	}
}
