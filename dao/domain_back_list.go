package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"mproxy/model"
)

func GetDomainBackListOnlyDomain(
	ctx context.Context,
	db *gorm.DB,
) ([]string, error) {
	var blacklist []string
	err := db.
		WithContext(ctx).Model(&model.IpRiskControlBlacklist{}).
		Where("is_deleted = ?", 0).
		Pluck("target_site", &blacklist).Error
	if err != nil {
		return nil, fmt.Errorf("推送域名黑名单到网关查询数据错误error:%+v", err)
	}

	return blacklist, nil
}
