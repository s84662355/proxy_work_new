package service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/model"
)

func BatchUpdateDynamicDataCache(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
) (RowsAffected int64, err error) {
	var results []*model.VsIPTransitDynamic = nil

	result := db.WithContext(ctx).Model(&model.VsIPTransitDynamic{}).FindInBatches(
		&results,
		constant.BatchUpdateDynamicDataCacheSize,
		func(tx *gorm.DB, batch int) error {
			for c := range slices.Chunk(results, 20) {
				if err := UpdateDynamicDataCachebyRedisPipe(
					ctx,
					rdb,
					c,
				); err != nil {
					log.Error("[service] BatchUpdateDynamicDataCache 设置数据 执行错误", zap.Any("error", err))
				}
			}

			return nil
		})
	return result.RowsAffected, result.Error
}

func UpdateDynamicDataCachebyRedisPipe(
	ctx context.Context,
	rdb *redis.Client,
	results []*model.VsIPTransitDynamic,
) error {
	_, err := rdb.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for _, v := range results {
				data, err := json.Marshal(v)
				if err != nil {
					log.Error("[service] UpdateDynamicDataCachebyRedisPipe json 解析失败", zap.Any("error", err), zap.Any("data", data))

					continue
				}
				s := string(data)

				pipe.Set(
					ctx,
					fmt.Sprintf("%s%d", constant.VsIPTransitDynamicCacheRedisKeyPrefix, v.UserID),
					s,
					constant.VsIPTransitDynamicCacheRedisTtl,
				)
			}
			return nil
		})

	if err == nil {
		return nil
	}

	return fmt.Errorf("UpdateDynamicDataCachebyRedisPipe %+v", err)
}
