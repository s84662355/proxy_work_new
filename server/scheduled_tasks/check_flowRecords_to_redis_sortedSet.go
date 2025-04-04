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

// 检查流量表
func (m *manager) checkFlowRecordsToRedisSortedSet(ctx context.Context) {
	loopTime := 60 * time.Second
	ticker := time.NewTicker(loopTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]检查流量表任务上下文结束退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks]定时执行检查流量表任务开始 ", zap.Time("Time", t))
			v, err := service.CheckFlowRecordsToRedisSortedSet(
				ctx,
				common.GetMysqlDB(),
				common.GetRedisDB(),
			)
			if err != nil {
				log.Error("[定时任务scheduled_tasks]定时执行检查流量表任务错误", zap.Any("error", err))
			}
			log.Info("[定时任务scheduled_tasks]定时执行检查流量表任务完成", zap.String("Time", fmt.Sprintf("%+v---%+v", t, time.Now())), zap.Int64("count", v))
			ticker.Reset(loopTime)
		}
	}
}

// 更新主账号流量
func (m *manager) updateFlowRecordsFromRedisSortedSet(ctx context.Context) {
	loopTime := 5 * time.Second
	ticker := time.NewTicker(loopTime)
	defer ticker.Stop()
	trm := taskRunManager.NewChanTask(5)
	for {
		ticker.Reset(loopTime)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]更新主账号流量任务上下文结束退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks]定时执行更新主账号流量开始 ", zap.Time("Time", t))
			_, results, err := common.GetRedisDB().
				BZMPop(
					ctx,
					5*time.Second, // 等待超时5秒
					"MIN",
					30, // 弹出30个
					constant.FlowUserIdQueueSortedSet,
				).Result()
			if err != nil {
				if err != redis.Nil {
					log.Error("[定时任务scheduled_tasks]更新主账号流量执行BZMPop命令时出错 ", zap.Any("error", err))
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
						log.Error("[定时任务scheduled_tasks]更新主账号流量执行strconv.ParseInt错误", zap.Any("error", err))
						continue
					}

					var recordId uint64 = 0
					if len(members) > 1 {
						recordId, _ = strconv.ParseUint(members[1], 10, 64)
					}

					if userId != 0 {
						trm.Run(func() {
							m.updateFlowAndAutoBuyFlow(userId, recordId)
						})
					}
				}
				trm.Wait()
			}

		}
	}
}

// 更新主账号流量并且主动购买流量
func (m *manager) updateFlowAndAutoBuyFlow(
	userId int64,
	recordId uint64,
) {
	///获取分布式锁
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
		log.Error("[定时任务scheduled_tasks]更新主账号流量并且主动购买流量错误", zap.Int64("userId", userId), zap.Any("error", err))
		return
	}

	///更新主账号缓存
	if err := service.UpdateDynamicDataCachebyRedisPipe(
		context.Background(),
		common.GetMysqlDB(),
		common.GetRedisDB(),
		userId,
	); err != nil {
		log.Error("[定时任务scheduled_tasks]更新主账号流量并且主动购买流量错误", zap.Int64("userId", userId), zap.Any("error", err))
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
	// 自动购买流量
	if _, err := service.DynamicFlowAndAutoBuy(
		context.Background(),
		common.GetMysqlDB(),
		userId,
	); err != nil {
		log.Error("[定时任务scheduled_tasks]自动购买流量错误", zap.Int64("userId", userId), zap.Any("error", err))
	}
}
