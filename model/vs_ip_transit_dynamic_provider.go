package model

const VsIPTransitDynamicProviderTableName = "vs_ip_transit_dynamic_provider"

// 供应商配置表
type VsIPTransitDynamicProvider struct {
	ProviderName string `gorm:"column:provider_name;type:varchar(50);primary_key" json:"provider_name"` // 供应商名称
	Proportion   uint   `gorm:"column:proportion;type:int(11) unsigned" json:"proportion"`              // 分配比例
	Locking      string `gorm:"column:locking;type:char(1)" json:"locking"`                             // 锁定 1-锁定 0-未锁定
	ResidualFlow int64  `gorm:"column:residual_flow;type:bigint(20)" json:"residual_flow"`              // 剩余流量（单位B）
	Setting      string `gorm:"column:setting;type:json" json:"setting"`                                // 供应商配置
	IsDatacenter bool   `gorm:"column:is_datacenter;type:json" json:"is_datacenter"`
}

func (m *VsIPTransitDynamicProvider) TableName() string {
	return VsIPTransitDynamicProviderTableName
}
