package dao

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// go test -v -run TestSelectPartitions  -tags "scheduled_tasks"
func TestSelectPartitions(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/new_schema?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	tableName := "vs_ip_flow_records"
	partitionName := "p20250328"
	exists, err := CheckPartitionExists(
		context.Background(),
		db,
		tableName,
		partitionName,
	)
	if err != nil {
		fmt.Println("Error checking partition existence:", err)
		return
	}
	if exists {
		fmt.Printf("Partition %s exists in table %s.\n", partitionName, tableName)
	} else {
		fmt.Printf("Partition %s does not exist in table %s.\n", partitionName, tableName)
	}
}

// go test -v -run TestCreatePartitions  -tags "scheduled_tasks"
func TestCreatePartitions(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/new_schema?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	// PARTITION p20250403 VALUES LESS THAN (1743695999) ENGINE = InnoDB)
	tableName := "vs_ip_flow_records"
	partitionName := "p20250403"
	err = CreateRangePartition(
		context.Background(),
		db,
		tableName,
		partitionName,
		1743695999,
	)
	if err != nil {
		fmt.Println("Error create partition error:", err)
		return
	}
}

// go test -v -run TestDeleteRangePartition  -tags "scheduled_tasks"
func TestDeleteRangePartition(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/new_schema?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	tableName := "vs_ip_flow_records"
	partitionName := "p20250403"
	err = DeleteRangePartition(
		context.Background(),
		db,
		tableName,
		partitionName,
	)
	if err != nil {
		fmt.Println("Error delete partition error:", err)
		return
	}
}
