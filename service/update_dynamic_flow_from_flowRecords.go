package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mproxy/constant"
	"mproxy/model"
)

func CheckFlowRecordsToRedisSortedSet(
	ctx context.Context,
	db *gorm.DB,
	rdb *redis.Client,
) (int64, error) {
	results := []*model.VsIPFlowRecords{}

	sql := fmt.Sprintf(`
           select user_id , MIN(id) as id
             from  
	           (  
	           		select 
	           			id,user_id 
	           				from ` + model.VsIPFlowRecordsTableName + ` 
	           					where is_deal = 0
	           						ORDER BY id ASC  
	           							LIMIT 500 
	           ) as Records 
	                GROUP BY user_id
    `)
	if err := db.WithContext(ctx).Raw(sql).Scan(&results).Error; err != nil {
		return 0, fmt.Errorf("CheckFlowRecordsToRedisSortedSet sql raw err:%+v", err)
	}

	var elements []redis.Z = nil
	for _, v := range results {
		elements = append(elements, redis.Z{
			Member: fmt.Sprint(v.UserID, ",", v.ID),
			Score:  float64(time.Now().Unix()),
		})
	}

	if len(elements) == 0 {
		return 0, nil
	}

	if v, err := rdb.ZAddNX(
		ctx,
		constant.FlowUserIdQueueSortedSet,
		elements...,
	).Result(); err != nil {
		return 0, fmt.Errorf("CheckFlowRecordsToRedisSortedSet ZAddNX err:%+v", err)
	} else {
		return v, nil
	}
}
