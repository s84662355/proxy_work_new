package model

import (
	"time"
)

const VsIPTransitDynamicWhitelistTableName = "vs_ip_transit_dynamic_whitelist"

// 白名单表
type VsIPTransitDynamicWhitelist struct {
	ID             int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"` // 白名单id
	IP             string    `gorm:"column:ip;type:varchar(50);NOT NULL" json:"ip"`                  // 白名单ip
	UserID         int64     `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`         // 创建用户id
	Remark         string    `gorm:"column:remark;type:varchar(255)" json:"remark"`                  // 备注
	LastAccessTime time.Time `gorm:"column:last_access_time;type:datetime" json:"last_access_time"`  // 最后访问时间
	CreateTime     time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`            // 创建时间
}

func (m *VsIPTransitDynamicWhitelist) TableName() string {
	return VsIPTransitDynamicWhitelistTableName
}
