package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mproxy/constant"
	"mproxy/model"
)

// /根据userid查询账号信息
func GetVsIPTransitDynamic(
	ctx context.Context,
	db *gorm.DB,
	userid uint64,
) (*model.VsIPTransitDynamic, error) {
	vsIPTransitDynamic := &model.VsIPTransitDynamic{}
	err := db.
		WithContext(ctx).
		Where("user_id = ?", userid).
		First(vsIPTransitDynamic).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("GetVsIPTransitDynamic查找user_id=%d数据不存在", userid)
		} else {
			return nil, fmt.Errorf("GetVsIPTransitDynamic查找user_id=%d数据错误 err:%+v", userid, err)
		}
	}
	return vsIPTransitDynamic, nil
}

// 设置主账号缓存
func SetVsIPTransitDynamicRedisCache(
	ctx context.Context,
	rdb *redis.Client,
	vsIPTransitDynamic *model.VsIPTransitDynamic,
) error {
	data, err := json.Marshal(vsIPTransitDynamic)
	if err != nil {
		return fmt.Errorf("SetVsIPTransitDynamicRedisCache  json.Marshal失败 数据:%+v err:%+v", vsIPTransitDynamic, err)
	}
	s := string(data)

	if _, err = rdb.Set(
		ctx,
		fmt.Sprintf("%s%d", constant.VsIPTransitDynamicCacheRedisKeyPrefix, vsIPTransitDynamic.UserID),
		s,
		constant.VsIPTransitDynamicCacheRedisTtl,
	).Result(); err != nil {
		return fmt.Errorf("SetVsIPTransitDynamicRedisCache  设置redis缓存失败 数据:%+v err:%+v", vsIPTransitDynamic, err)
	}

	return nil
}
