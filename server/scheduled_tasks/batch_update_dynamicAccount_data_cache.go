package scheduled_tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/service"
)

// /更新子账号redis缓存
func (m *manager) batchUpdateDynamicAccountDataCache(ctx context.Context) {
	if c, err := service.BatchUpdateDynamicAccountDataCache(
		ctx,
		common.GetMysqlDB(),
		common.GetRedisDB(),
	); err != nil {
		log.Error("[定时任务scheduled_tasks] 首次执行 batchUpdateDynamicAccountDataCache 执行错误", zap.Int64("count", c), zap.Any("error", err))
	}
	ticker := time.NewTicker(constant.DynamicAccountBatchUpdateDataLoopTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks] 定时执行 batchUpdateDynamicAccountDataCache 上下文Done() 退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks] 定时执行 batchUpdateDynamicAccountDataCache 开始", zap.Any("Time", t))
			c, err := service.BatchUpdateDynamicAccountDataCache(
				ctx,
				common.GetMysqlDB(),
				common.GetRedisDB(),
			)
			if err != nil {
				log.Error("[定时任务scheduled_tasks] 定时执行 batchUpdateDynamicAccountDataCache", zap.Any("Time", t), zap.Int64("count", c), zap.Any("error", err))
			}
			log.Info("[定时任务scheduled_tasks] 定时执行 batchUpdateDynamicAccountDataCache 结束", zap.String("Time", fmt.Sprintf("%+v---%+v", t, time.Now())), zap.Int64("count", c))
			ticker.Reset(constant.DynamicAccountBatchUpdateDataLoopTime)
		}
	}
}
