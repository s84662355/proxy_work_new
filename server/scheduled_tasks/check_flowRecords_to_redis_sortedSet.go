package scheduled_tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/log"
	"mproxy/service"
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
