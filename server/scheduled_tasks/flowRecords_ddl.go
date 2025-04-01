package scheduled_tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/log"
	"mproxy/model"
)

// 删除流量表分区和新增分区
func (m *manager) updateVsIpFlowRecordsPartition(
	ctx context.Context,
) {
	ticker, done := m.updateVsIpFlowRecordsPartitionAction()
	defer func() {
		ticker.Stop()
		for range done {
		}
		log.Info("[定时任务scheduled_tasks]删除流量表分区和新增分区任务完全退出")
	}()

	log.Info("[定时任务scheduled_tasks]删除流量表分区和新增分区任务开始")
	<-ctx.Done()
	log.Info("[定时任务scheduled_tasks]删除流量表分区和新增分区任务上下文结束退出")
}

// 时执行删除流量表分区和新增分区
func (m *manager) updateVsIpFlowRecordsPartitionAction() (*time.Ticker, <-chan struct{}) {
	ticker := time.NewTicker(CalculationDifference())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range ticker.C {
			log.Info("[定时任务scheduled_tasks]定时执行删除流量表分区和新增分区开始")
			// 获取当前时间
			now := time.Now()

			// 计算一周前的时间
			oneWeekAgo := now.AddDate(0, 0, -7)

			// 计算当天最后时间点的Unix时间戳
			todayEnd := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				23, 59, 59, 0,
				now.Location(),
			).Unix()

			// 构建删除分区的SQL
			dropPartitionSQL := fmt.Sprintf(
				"ALTER TABLE %s DROP PARTITION p%s",
				model.VsIPFlowRecordsTableName,
				oneWeekAgo.Format("20060102"),
			)

			// 构建新增分区的SQL
			addPartitionSQL := fmt.Sprintf(
				"ALTER TABLE %s ADD PARTITION (PARTITION p%s VALUES LESS THAN (%d))",
				model.VsIPFlowRecordsTableName,
				now.Format("20060102"),
				todayEnd,
			)

			// 执行删除分区
			if err := common.GetMysqlDB().Exec(dropPartitionSQL).Error; err != nil {
				log.Error("[scheduled_tasks]定时执行删除流量表分区和新增分区 删除分区失败", zap.Any("error", err))
			}

			// 执行新增分区
			if err := common.GetMysqlDB().Exec(addPartitionSQL).Error; err != nil {
				log.Error("[scheduled_tasks]定时执行删除流量表分区和新增分区 新增分区失败", zap.Any("error", err))
			}

			ticker.Reset(CalculationDifference())

		}
	}()

	return ticker, done
}

// 计算下一个0点时间
func CalculationDifference() time.Duration {
	// 计算下一个0点时间
	now := time.Now()
	nextMidnight := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	).Add(24 * time.Hour)
	waitDuration := nextMidnight.Sub(now)

	return waitDuration
}
