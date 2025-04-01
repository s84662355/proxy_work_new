package service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/model"
)

// 批量更新子账号缓存
func BatchUpdateDynamicAccountDataCache(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
) (RowsAffected int64, err error) {
	var results []*model.VsIPTransitDynamicAccount = nil
	result := db.WithContext(ctx).Model(&model.VsIPTransitDynamicAccount{}).FindInBatches(
		&results,
		constant.BatcheUpdateDynamicAccountDataCacheSize,
		func(tx *gorm.DB, batch int) error {
			for c := range slices.Chunk(results, 20) {
				///设置 DynamicAccount 缓存
				if err := UpdateDynamicAccountDataCachebyRedisPipe(ctx, rdb, c); err != nil {
					log.Error("[service]批量更新子账号缓存 设置数据执行错误", zap.Any("error", err))
				}

				if err := ExistsFlowDynamicAccountIDbyRedisPipe(ctx, rdb, results); err != nil {
					log.Error("[service]批量更新子账号缓存 Exists执行错误 ", zap.Any("error", err))
				}

			}

			return nil
		})
	return result.RowsAffected, result.Error
}

// 使用redis管道批量设置子账号缓存
func UpdateDynamicAccountDataCachebyRedisPipe(
	ctx context.Context,
	rdb *redis.Client,
	results []*model.VsIPTransitDynamicAccount,
) error {
	_, err := rdb.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for _, v := range results {
				if v.IsDelete == "0" {
					accountData, err := json.Marshal(v)
					if err != nil {
						log.Error("[service]使用redis管道批量设置子账号缓存 json解析失败", zap.Any("account", v), zap.Any("error", err))
						continue
					}
					s := string(accountData)
					pipe.Set(
						ctx,
						constant.DynamicAccountDataCacheRedisKeyPrefix+v.Username,
						s,
						constant.DynamicAccountDataCacheRedisTtl,
					)
					pipe.Set(
						ctx,
						fmt.Sprintf("%s%d", constant.DynamicAccountDataCacheByIdRedisKeyPrefix, v.ID),
						s,
						constant.DynamicAccountDataCacheRedisTtl,
					)

				}
			}
			return nil
		})

	if err == nil {
		return nil
	}

	return fmt.Errorf("使用redis管道批量设置子账号缓存 error:%+v", err)
}

// 使用redis管道批量判断子账号的流量是否存在
func ExistsFlowDynamicAccountIDbyRedisPipe(
	ctx context.Context,
	rdb *redis.Client,
	results []*model.VsIPTransitDynamicAccount,
) error {
	///后面可以试试使用MGet一次性获取所有key的值
	incrByCmds := map[int64]*redis.IntCmd{}
	_, err := rdb.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for _, v := range results {
				incrByCmd := pipe.IncrBy(
					ctx,
					fmt.Sprintf("%s%d", constant.DynamicAccountRedisFlowPrefix, v.ID),
					0,
				)
				incrByCmds[v.ID] = incrByCmd
			}
			return nil
		})
	if err != nil {
		err = fmt.Errorf("使用redis管道批量判断子账号的流量是否存在 redis Pipelined Exists执行错误 error:%+v", err)
		return err
	}

	var elements []redis.Z = nil
	///判断key是否存在并且值大于0
	for k, incrByCmd := range incrByCmds {
		result, err := incrByCmd.Result()
		if err == nil && result > 0 {
			elements = append(elements, redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: k,
			})
		}
	}

	if len(elements) == 0 {
		return nil
	}

	if _, err := rdb.ZAddNX(
		ctx,
		constant.DynamicAccountIDFlowRedisQueueSortedSet,
		elements...,
	).Result(); err != nil {
		return fmt.Errorf("使用redis管道批量判断子账号的流量是否存在 ZAddNX error:%+v", err)
	}

	return nil
}
