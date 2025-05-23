// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
	"time"
)

const VsIPTransitDynamicAccountAccessRecordsTableName = "vs_ip_transit_dynamic_account_access_records"

// 子账号访问记录表
type VsIPTransitDynamicAccountAccessRecords struct {
	ID         uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"` // 主键id
	UserID     int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`                           // 用户id
	AccountID  int64     `gorm:"column:account_id;type:bigint(20)" json:"account_id"`                     // 子账号id
	Username   string    `gorm:"column:username;type:varchar(200)" json:"username"`                       // 子账号用户名
	Domain     string    `gorm:"column:domain;type:varchar(200)" json:"domain"`                           // 访问域名
	CreateTime int64     `gorm:"column:create_time;type:bigint(20)" json:"create_time"`                   // 创建时间(时间戳)
	PoolType   string    `gorm:"column:pool_type;type:varchar(100)" json:"pool_type"`                     // 供应商名
	ProviderIP string    `gorm:"column:provider_ip;type:varchar(50)" json:"provider_ip"`                  // 供应商ip
	AccessTime time.Time `gorm:"column:access_time;type:datetime" json:"access_time"`                     // 时间
	IsConnect  string    `gorm:"column:is_connect;type:varchar(2)" json:"is_connect"`                     // 是否连通（0: 否 1：是）
}

func (m *VsIPTransitDynamicAccountAccessRecords) TableName() string {
	return VsIPTransitDynamicAccountAccessRecordsTableName
}
