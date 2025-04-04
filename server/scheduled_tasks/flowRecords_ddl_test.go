package scheduled_tasks

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// go test -v -run TestUpdateVsIpFlowRecordsPartition  -tags "scheduled_tasks"
func TestUpdateVsIpFlowRecordsPartition(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/new_schema?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	m := &manager{}

	m.updateVsIpFlowRecordsPartitionAction(db)
}
