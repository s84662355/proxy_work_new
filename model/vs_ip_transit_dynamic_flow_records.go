package model

import (
	"time"
)

type VsIPTransitDynamicFlowRecords struct {
	ID                        int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	OldUseFlow                uint64    `gorm:"column:old_use_flow;type:bigint(20) unsigned;NOT NULL" json:"old_use_flow"`                       // 更新前的使用流量
	UseFlow                   uint64    `gorm:"column:use_flow;type:bigint(20) unsigned;NOT NULL" json:"use_flow"`                               // 更新后的使用流量
	OldUseFlowDatacenter      uint64    `gorm:"column:old_use_flow_datacenter;type:bigint(20) unsigned;NOT NULL" json:"old_use_flow_datacenter"` // 更新前的机房使用流量
	UseFlowDatacenter         uint64    `gorm:"column:use_flow_datacenter;type:bigint(20) unsigned;NOT NULL" json:"use_flow_datacenter"`         // 更新后的机房使用流量
	Flow                      int64     `gorm:"column:flow;type:bigint(20);NOT NULL" json:"flow"`                                                // 本次更新的流量
	FlowDatacenter            int64     `gorm:"column:flow_datacenter;type:bigint(20);NOT NULL" json:"flow_datacenter"`                          // 本次更新的流量
	UserID                    int64     `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`                                          // 系统账号id
	FlowRecordsIds            string    `gorm:"column:flow_records_ids;type:text;NOT NULL" json:"flow_records_ids"`                              // 记录id，使用逗号隔开
	CreateTime                time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
	OldResidualFlow           int64     `gorm:"column:old_ residual_flow;type:bigint(20);NOT NULL" json:"old_ residual_flow"`                     // 更新前的剩余流量
	ResidualFlow              int64     `gorm:"column:residual_flow;type:bigint(20);NOT NULL" json:"residual_flow"`                               // 更新后的剩余流量
	OldResidualFlowDatacenter int64     `gorm:"column:old_residual_flow_datacenter;type:bigint(20);NOT NULL" json:"old_residual_flow_datacenter"` // 更新前的剩余流量
	ResidualFlowDatacenter    int64     `gorm:"column:residual_flow_datacenter;type:bigint(20);NOT NULL" json:"residual_flow_datacenter"`         // 更新后的剩余流量
}

func (m VsIPTransitDynamicFlowRecords) TableName() string {
	return "vs_ip_transit_dynamic_flow_records"
}
