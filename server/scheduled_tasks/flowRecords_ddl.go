package scheduled_tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"mproxy/common"
	"mproxy/log"
	"mproxy/model"
)

func (m *manager) updateVsIpFlowRecordsPartition(
	ctx context.Context,
) {
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
	ticker := time.NewTicker(waitDuration)

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks] 定时执行 updateVsIpFlowRecordsPartition 上下文Done() 退出")
			return
		case t := <-ticker.C:
			log.Info("[定时任务scheduled_tasks] 定时执行 updateVsIpFlowRecordsPartition 开始 ", zap.Time("Time", t))
			m.updateVsIpFlowRecordsPartitionAction(
				common.GetMysqlDB(),
			)

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
			ticker.Reset(waitDuration)
		}
	}
}

func (m *manager) updateVsIpFlowRecordsPartitionAction(
	db *gorm.DB,
) {
	// 获取当前时间
	now := time.Now()

	// 计算一周前的时间
	oneWeekAgo := now.AddDate(0, 0, -7)

	// 计算当天最后时间点的Unix时间戳
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Unix()

	// 构建删除分区的SQL
	dropPartitionSQL := fmt.Sprintf("ALTER TABLE %s DROP PARTITION p%s", model.VsIPFlowRecordsTableName,
		oneWeekAgo.Format("20060102"))

	// 构建新增分区的SQL
	addPartitionSQL := fmt.Sprintf(
		"ALTER TABLE %s ADD PARTITION (PARTITION p%s VALUES LESS THAN (%d))",
		model.VsIPFlowRecordsTableName,
		now.Format("20060102"),
		todayEnd,
	)

	// 执行删除分区
	if err := db.Exec(dropPartitionSQL).Error; err != nil {
		log.Error("[scheduled_tasks] UpdateVsIpFlowRecordsPartition 删除分区失败", zap.Any("error", err))
	}

	// 执行新增分区
	if err := db.Exec(addPartitionSQL).Error; err != nil {
		log.Error("[scheduled_tasks] UpdateVsIpFlowRecordsPartition 新增分区失败", zap.Any("error", err))
	}
}
