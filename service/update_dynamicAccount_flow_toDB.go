package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mproxy/constant"
	"mproxy/dao"
	"mproxy/log"
	"mproxy/model"
)

func UpdateDynamicAccountFlowFromRedisToDBWithRedisLock() {
}

// /从redis更新子账号流量
func UpdateDynamicAccountFlowFromRedisToDB(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
	accountId int64,
) error {
	IncrementKey := fmt.Sprintf("%s%d", constant.DynamicAccountRedisFlowPrefix, accountId)
	flow, err := rdb.GetDel(ctx, IncrementKey).Result()
	if err != nil {
		return fmt.Errorf("从redis更新子账号流量 rdb.GetDel err:%+v", err)
	}

	supplyFlow, err := strconv.ParseInt(flow, 10, 64)
	if err != nil {
		return fmt.Errorf("从redis更新子账号流量  strconv.ParseInt  err:%+v", err)
	}

	if supplyFlow == 0 {
		return nil
	}

	err = UpdateDynamicAccountFlowToDB(
		ctx,
		db,
		accountId,
		supplyFlow,
	)
	if err != nil {
		if er := dao.RedisIncrByExpirationOnPipe(
			context.Background(),
			rdb,
			IncrementKey,
			supplyFlow,
			constant.DynamicAccountRedisFlowTtl,
		); er != nil {
			log.Error("[service]从redis更新子账号流量 回滚流量到redis失败", zap.Any("error", err))
		}

		return err
	}

	return nil
}

// /更新子账号的流量并且插入流量记录
func UpdateDynamicAccountFlowToDB(
	ctx context.Context,
	db *gorm.DB,
	accountId int64,
	supplyFlow int64,
) error {
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		////修改账号的流量
		errTransaction := tx.Exec(
			fmt.Sprintf(
				`update %s SET use_flow = use_flow + ? , update_time = ? WHERE id = ?`,
				model.VsIPTransitDynamicAccountTableName,
			),
			supplyFlow,
			time.Now(),
			accountId,
		).Error
		if errTransaction != nil {
			return errTransaction
		}

		// 更新子账号每日的流量
		flowSQl := `
				REPLACE into ` + model.VsIPTransitDynamicAccountFlowTableName + ` 
					(account_id ,date_time,use_flow)  
						(
							select  
								? as account_id, 
								DATE_FORMAT(NOW(),'%Y%m%d') as date_time, 
								ifnull
								(
									( 
										SELECT  
											(use_flow + ?)  as use_flow  
												from  ` + model.VsIPTransitDynamicAccountFlowTableName + `
													 where 
													 	account_id = ?  
													 		and 	
													 	date_time = DATE_FORMAT(NOW(),'%Y%m%d') 
									),
								 	?
								) as use_flow 
						)
			`

		////记录子账号当天使用的流量
		return tx.Exec(
			flowSQl,
			accountId,
			supplyFlow,
			accountId,
			supplyFlow,
		).Error
	}); err != nil {
		return fmt.Errorf("更新子账号的流量 并且插入流量记录  db.Transaction err:%+v", err)
	}

	return nil
}
