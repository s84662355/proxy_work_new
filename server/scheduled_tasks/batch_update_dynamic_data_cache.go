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

// 更新主账号redis缓存
func (m *manager) batchUpdateDynamicDataCache(ctx context.Context) {
	if c, err := service.BatchUpdateDynamicDataCache(
		ctx,
		common.GetMysqlDB(),
		common.GetRedisDB(),
	); err != nil {
		log.Error("[定时任务scheduled_tasks]首次执行更新主账号redis缓存错误", zap.Int64("count", c), zap.Any("error", err))
	}
	ticker := time.NewTicker(constant.VsIPTransitDynamicBatchUpdateDataLoopTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]更新主账号redis缓存任务上下文结束退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks]定时执行更新主账号redis缓存开始")
			c, err := service.BatchUpdateDynamicDataCache(
				ctx,
				common.GetMysqlDB(),
				common.GetRedisDB(),
			)
			if err != nil {
				log.Error("[定时任务scheduled_tasks]定时执行更新主账号redis缓存", zap.Int64("count", c), zap.Any("error", err))
			}
			log.Info("[定时任务scheduled_tasks]]定时执行更新主账号redis缓存完成", zap.String("Time", fmt.Sprintf("%+v---%+v", t, time.Now())), zap.Int64("count", c))
			ticker.Reset(constant.VsIPTransitDynamicBatchUpdateDataLoopTime)
		}
	}
}
