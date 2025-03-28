package model

import (
	_ "fmt"
	"time"
)

type IpRiskControlAccessLog struct {
	ID               uint64    `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	TargetSite       string    `gorm:"column:target_site;type:varchar(255);NOT NULL" json:"target_site"`                 // 地址
	AccountType      int       `gorm:"column:account_type;type:tinyint(4);NOT NULL" json:"account_type"`                 // 账号类型 0动态 1静态
	Account          string    `gorm:"column:account;type:varchar(255);NOT NULL" json:"account"`                         // 账户
	ExitIp           string    `gorm:"column:exit_ip;type:varchar(45);NOT NULL" json:"exit_ip"`                          // 出口ip
	AccessCount      int       `gorm:"column:access_count;type:int(10) unsigned;default:0;NOT NULL" json:"access_count"` // 访问次数
	LatestAccessTime time.Time `gorm:"column:latest_access_time;type:datetime" json:"latest_access_time"`                // 最新访问时间
	CreateAt         time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`                                // 创建的时间
}

func (m IpRiskControlAccessLog) TableName() string {
	return "ip_risk_control_access_log"
}
