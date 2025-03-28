package model

// 子账号每日使用流量表
type VsIPTransitDynamicAccountFlow struct {
	AccountID         uint64 `gorm:"primaryKey;autoIncrement:false;column:account_id;type:bigint(20) unsigned;NOT NULL" json:"account_id"` // 子账号id
	DateTime          uint64 `gorm:"primaryKey;autoIncrement:false;column:date_time;type:bigint(10) unsigned;NOT NULL" json:"date_time"`   // 流量使用时间
	UseFlow           uint64 `gorm:"column:use_flow;type:bigint(20) unsigned;NOT NULL" json:"use_flow"`                                    // 已使用流量
	UseFlowDatacenter uint64 `gorm:"column:use_flow_datacenter;type:bigint(20) unsigned;NOT NULL" json:"use_flow_datacenter"`              // 已使用静态机房流量
}

func (m *VsIPTransitDynamicAccountFlow) TableName() string {
	return "vs_ip_transit_dynamic_account_flow"
}
