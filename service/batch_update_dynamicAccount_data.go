package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mproxy/common"
	"mproxy/constant"
	"mproxy/log"
	"mproxy/model"
)

// /批量更新子账号缓存
func BatchUpdateDynamicAccountData(ctx context.Context) (RowsAffected int64, err error) {
	var (
		gormDB                                     = common.GetMysqlDB()
		results []*model.VsIPTransitDynamicAccount = nil
	)

	// 处理记录，批处理大小为100
	result := gormDB.Model(&model.VsIPTransitDynamicAccount{}).FindInBatches(&results, constant.BatcheUpdateDynamicAccountDataCacheSize, func(tx *gorm.DB, batch int) error {
		///设置 DynamicAccount 缓存
		if err := UpdateDynamicAccountDatabyRedisPipe(ctx, results); err != nil {
			log.Errorf("[service] BatchUpdateDynamicAccountTraffic 设置数据 执行错误 err:%+v", err)
		}

		aids, err := ExistsFlowDynamicAccountIDbyRedisPipe(ctx, results)
		if err != nil {
			log.Errorf("[service] BatchUpdateDynamicAccountTraffic Exists   执行错误 err:%+v", err)
		}

		if len(aids) > 0 {
			// 批量添加成员到集合
			_, err := common.GetRedisDB().SAdd(ctx, constant.DynamicAccountIDFlowRedisQueueSet, aids...).Result()
			if err != nil {
				log.Errorf("[service] SAddHaveFlowDynamicAccountIDToQueueSet redis SAdd %s 执行错误 err:%+v", constant.DynamicAccountIDFlowRedisQueueSet, err)
			}
		}

		return nil
	})
	return result.RowsAffected, result.Error
}

func UpdateDynamicAccountDatabyRedisPipe(ctx context.Context, results []*model.VsIPTransitDynamicAccount) error {
	_, err := common.GetRedisDB().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, v := range results {
			if v.IsDelete == "0" {
				accountData, err := json.Marshal(v)
				if err != nil {
					log.Errorf("[service] UpdateDynamicAccountDataOnRedisPipe json 解析失败 数据%+v err:%+v", v, err)
					continue
				}
				pipe.Set(ctx, constant.DynamicAccountDataCacheRedisKeyPrefix+v.Username, string(accountData), constant.DynamicAccountDataCacheRedisTtl)
				pipe.Set(ctx, fmt.Sprintf("%s%d", constant.DynamicAccountDataCacheByIdRedisKeyPrefix, v.ID), string(accountData), constant.DynamicAccountDataCacheRedisTtl)

			}
		}

		return nil
	})

	if err == nil {
		return nil
	}

	return fmt.Errorf("UpdateDynamicAccountDataOnRedisPipe %+v", err)
}

func ExistsFlowDynamicAccountIDbyRedisPipe(ctx context.Context, results []*model.VsIPTransitDynamicAccount) ([]interface{}, error) {
	existsCmds := map[int64]*redis.IntCmd{}
	_, err := common.GetRedisDB().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, v := range results {
			existsCmd := pipe.Exists(ctx, fmt.Sprintf("%s%d", constant.DynamicAccountRedisFlowPrefix, v.ID))
			existsCmds[v.ID] = existsCmd
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("ExistsFlowDynamicAccountIDbyRedisPipe redis Pipelined Exists 执行错误 err:%+v", err)
	}

	var aids []interface{} = nil
	for k, existsCmd := range existsCmds {
		result, err := existsCmd.Result()
		if err == nil && result > 0 {
			aids = append(aids, k)
		}
	}

	return aids, err
}
