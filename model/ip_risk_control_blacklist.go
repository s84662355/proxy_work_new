package model

import "time"

type IpRiskControlBlacklist struct {
	ID               uint64    `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	TargetSite       string    `gorm:"column:target_site;type:varchar(255);NOT NULL" json:"target_site"`          // 地址
	IsDeleted        int       `gorm:"column:is_deleted;type:tinyint(4);default:0;NOT NULL" json:"is_deleted"`    // 是否删除 0未删除 1已删除
	CreateAt         time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`                         // 创建的时间
	AlarmThreshold   float64   `gorm:"column:alarm_threshold;type:decimal(10,2);NOT NULL" json:"alarm_threshold"` // 告警阈值
	LatestAccessTime time.Time `gorm:"column:latest_access_time;type:datetime" json:"latest_access_time"`         // 最新访问时间
}

func (m IpRiskControlBlacklist) TableName() string {
	return "ip_risk_control_blacklist"
}
