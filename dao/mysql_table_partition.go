package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// /判断分区是否存在
func CheckPartitionExists(
	ctx context.Context,
	db *gorm.DB,
	tableName string,
	partitionName string,
) (bool, error) {
	var count int64
	err := db.
		WithContext(ctx).
		Raw(
			`
					SELECT COUNT(*) 
						FROM information_schema.partitions 
							WHERE table_schema = DATABASE()
								AND table_name =? AND 
									partition_name =?
				`,
			tableName,
			partitionName,
		).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 创建分区
func CreateRangePartition(
	ctx context.Context,
	db *gorm.DB,
	tableName string,
	partitionName string,
	intervalValue int64,
) error {
	// 构建新增分区的SQL
	addPartitionSQL := fmt.Sprintf(
		`
			ALTER 
				TABLE %s 
					ADD PARTITION 
						(
							PARTITION %s VALUES LESS THAN (%d)
						)
		`,
		tableName,
		partitionName,
		intervalValue,
	)

	// 执行新增分区
	if err := db.
		WithContext(ctx).
		Exec(addPartitionSQL).
		Error; err != nil {
		return fmt.Errorf("CreateRangePartition创建分区失败error:%+v", err)
	}

	return nil
}

// /删除分区
func DeleteRangePartition(
	ctx context.Context,
	db *gorm.DB,
	tableName string,
	partitionName string,
) error {
	// 构建删除分区的SQL
	dropPartitionSQL := fmt.Sprintf(
		`
			ALTER 
				TABLE %s 
					DROP 
						PARTITION %s
		`,
		tableName,
		partitionName,
	)

	// 执行删除分区
	if err := db.
		WithContext(ctx).
		Exec(dropPartitionSQL).
		Error; err != nil {
		return fmt.Errorf("DeleteRangePartition删除分区失败:%+v", err)
	}

	return nil
}
