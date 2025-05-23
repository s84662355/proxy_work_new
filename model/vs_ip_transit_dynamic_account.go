// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
	"encoding/json"
	"time"
)

const VsIPTransitDynamicAccountTableName = "vs_ip_transit_dynamic_account"

// 动态ip子用户表
type VsIPTransitDynamicAccount struct {
	ID                   int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`                     // 主键id
	UserID               int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`                                      // 系统用户id
	ParentID             int64     `gorm:"column:parent_id;type:bigint(20);default:-1" json:"parent_id"`                       // 父id（主账号为-1）
	Status               string    `gorm:"column:status;type:char(1);default:0" json:"status"`                                 // 状态（1-启用中 0-未启用）
	Username             string    `gorm:"column:username;type:varchar(100)" json:"username"`                                  // 用户名（用户账号）
	Password             string    `gorm:"column:password;type:varchar(100)" json:"password"`                                  // 密码
	UseFlow              uint64    `gorm:"column:use_flow;type:bigint(20) unsigned" json:"use_flow"`                           // 已使用流量（单位B）
	UpperFlow            uint64    `gorm:"column:upper_flow;type:bigint(20) unsigned" json:"upper_flow"`                       // 上限流量（单位B 0为没有限制）
	UseFlowDatacenter    uint64    `gorm:"column:use_flow_datacenter;type:bigint(20) unsigned" json:"use_flow_datacenter"`     // 已使用流量机房（单位B）
	UpperFlowDatacenter  uint64    `gorm:"column:upper_flow_datacenter;type:bigint(20) unsigned" json:"upper_flow_datacenter"` // 上限流量机房（单位B 0为没有限制）
	Remark               string    `gorm:"column:remark;type:varchar(255)" json:"remark"`                                      // 备注
	CreateTime           time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`                                // 创建时间
	UpdateTime           time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`                                // 修改时间
	IsDelete             string    `gorm:"column:is_delete;type:char(1);default:0" json:"is_delete"`                           // 是否删除（1-是 0-否）
	DirectConnectionInfo string    `gorm:"column:direct_connection_info;type:varchar(100)" json:"direct_connection_info"`      // 直连信息（国家城市逗号隔开）
}

func (m VsIPTransitDynamicAccount) TableName() string {
	return VsIPTransitDynamicAccountTableName
}

type DirectConnectionInfo struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Time    int64  `json:"time"`
}

func (m *VsIPTransitDynamicAccount) GetDirectConnectionInfo() *DirectConnectionInfo {

	directConnectionInfo := &DirectConnectionInfo{}
	if m.DirectConnectionInfo == "" {
		return directConnectionInfo
	}

	err := json.Unmarshal([]byte(m.DirectConnectionInfo), directConnectionInfo)
	if err != nil {
		return nil
	}

	return directConnectionInfo
}
