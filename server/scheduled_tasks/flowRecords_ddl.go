package scheduled_tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"mproxy/common"
	"mproxy/dao"
	"mproxy/log"
	"mproxy/model"
)

// 删除流量表分区和新增分区
func (m *manager) updateVsIpFlowRecordsPartition(
	ctx context.Context,
) {
	loopTime := 15 * 60 * time.Second
	ticker := time.NewTicker(loopTime)
	for {
		ticker.Reset(loopTime)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]删除流量表分区和新增分区任务上下文结束退出")
			return
		case <-ticker.C:
			log.Info("[定时任务scheduled_tasks]定时执行删除流量表分区和新增分区开始")
			m.updateVsIpFlowRecordsPartitionAction(
				common.GetMysqlDB(),
			)
		}
	}
}

// 时执行删除流量表分区和新增分区
func (m *manager) updateVsIpFlowRecordsPartitionAction(
	db *gorm.DB,
) {
	// 获取现在的时间
	timeDate := time.Now()

	{
		// 创建今天的分区
		m.createFlowRecordsPartition(
			timeDate,
			db,
		)
	}

	{
		// 创建明天的分区
		timeDate = timeDate.AddDate(0, 0, 1)
		m.createFlowRecordsPartition(
			timeDate,
			db,
		)
	}

	{
		// 创建后天的分区
		timeDate = timeDate.AddDate(0, 0, 1)
		m.createFlowRecordsPartition(
			timeDate,
			db,
		)
	}

	{
		// 删除7天前的分区
		m.deleteFlowRecordsPartition(
			time.Now().AddDate(0, 0, -7),
			db,
		)
	}
}

// 创建分区
func (m *manager) createFlowRecordsPartition(
	t time.Time,
	db *gorm.DB,
) {
	// 按照 "20060102" 格式格式化当前时间，得到今天的日期
	partitionName := fmt.Sprintf("p%s", t.Format("20060102"))
	exists, err := dao.CheckPartitionExists(
		context.Background(),
		db,
		model.VsIPFlowRecordsTableName,
		partitionName,
	)
	if err != nil {
		log.Error("[scheduled_tasks]定时执行新增分区 检查分区失败", zap.String("partition", partitionName), zap.Any("error", err))
	}

	if !exists {
		year, month, day := t.Date()
		//  23:59:59 的时间
		lastSecondOfToday := time.Date(year, month, day, 23, 59, 59, 0, t.Location())
		err := dao.CreateRangePartition(
			context.Background(),
			db,
			model.VsIPFlowRecordsTableName,
			partitionName,
			lastSecondOfToday.Unix(),
		)
		if err != nil {
			log.Error("[scheduled_tasks]定时执行新增分区失败", zap.String("partition", partitionName), zap.Any("error", err))
			return
		}
		log.Info("[定时任务scheduled_tasks]定时执行创建分区成功", zap.String("partition", partitionName))
	}
}

func (m *manager) deleteFlowRecordsPartition(
	t time.Time,
	db *gorm.DB,
) {
	// 按照 "20060102" 格式格式化当前时间，得到今天的日期
	partitionName := fmt.Sprintf("p%s", t.Format("20060102"))
	exists, err := dao.CheckPartitionExists(
		context.Background(),
		db,
		model.VsIPFlowRecordsTableName,
		partitionName,
	)
	if err != nil {
		log.Error("[scheduled_tasks]定时执行删除分区 检查分区失败", zap.String("partition", partitionName), zap.Any("error", err))
	}

	if exists {
		err := dao.DeleteRangePartition(
			context.Background(),
			db,
			model.VsIPFlowRecordsTableName,
			partitionName,
		)
		if err != nil {
			log.Error("[scheduled_tasks]定时执行删除分区失败", zap.String("partition", partitionName), zap.Any("error", err))
			return
		}
		log.Info("[定时任务scheduled_tasks]定时执行删除分区成功", zap.String("partition", partitionName))
	}
}
