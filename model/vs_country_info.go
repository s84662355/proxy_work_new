package model

import (
	"time"
)

const VsCountryInfoTableName = "vs_country_info"

// 虚拟服务国家信息表
type VsCountryInfo struct {
	ID          uint       `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(64);NOT NULL" json:"name"`                      // 国家名称
	Desc        string     `gorm:"column:desc;type:varchar(64);NOT NULL" json:"desc"`                      // 描述
	ShortName   string     `gorm:"column:short_name;type:char(8);NOT NULL" json:"short_name"`              // 简写名称
	PhoneCode   uint       `gorm:"column:phone_code;type:smallint(6) unsigned;NOT NULL" json:"phone_code"` // 区号
	Timelag     string     `gorm:"column:timelag;type:char(5);NOT NULL" json:"timelag"`                    // 时差
	Banner      *string    `gorm:"column:banner;type:varchar(255)" json:"banner"`                          // 国旗
	GmtCreate   *time.Time `gorm:"column:gmt_create;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_create"`
	GmtModified *time.Time `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"`
}

func (m *VsCountryInfo) TableName() string {
	return VsCountryInfoTableName
}
