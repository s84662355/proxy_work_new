package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mproxy/api"
	_ "mproxy/constant"
	"mproxy/dao"
	"mproxy/model"
)

func DynamicFlowAndAutoBuy(
	ctx context.Context,
	db *gorm.DB,
	userId int64,
) (bool, error) {
	vsIPTransitDynamic, err := dao.GetVsIPTransitDynamic(ctx, db, uint64(userId))
	if err != nil {
		return false, fmt.Errorf("UpdateFlowRecordsToDynamicFlowAndAutoBuy 获取用户信息失败 err:%+v", err)
	}
	normalFlowAutoBuy := false
	datacenterFlowAutoBuy := false
	if vsIPTransitDynamic.AutoBuy == "1" && vsIPTransitDynamic.ResidualFlow <= vsIPTransitDynamic.LastBuy {
		normalFlowAutoBuy = true
	}
	if vsIPTransitDynamic.AutoBuyDatacenter == "1" && vsIPTransitDynamic.ResidualFlowDatacenter <= vsIPTransitDynamic.LastBuyDatacenter {
		datacenterFlowAutoBuy = true
	}
	if normalFlowAutoBuy || datacenterFlowAutoBuy {
		///调用自动购买流量接口
		_, err := api.AutomaticRechargeTraffic(userId, normalFlowAutoBuy, datacenterFlowAutoBuy)
		if err != nil {
			return true, fmt.Errorf(
				`UpdateFlowRecordsToDynamicFlowAndAutoBuy normalFlowAutoBuy:%+v, datacenterFlowAutoBuy:%+v  自动购买流量失败 err:%+v`,
				normalFlowAutoBuy,
				datacenterFlowAutoBuy,
				err,
			)
		}

		return true, nil
	}

	return false, nil
}

// /更新流量数据记录到主账号
func UpdateFlowRecordsToDynamicFlow(
	ctx context.Context,
	db *gorm.DB,
	userId int64,
	flowIncRate float64,
	limit int,
	recordId uint64,
) (flow, flowDatacenter int64, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		vsIPTransitDynamic := &model.VsIPTransitDynamic{}
		err = tx.Where("user_id = ?  ", userId).First(vsIPTransitDynamic).Error
		if err != nil {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找id:%d主账号信息错误 err:%+v", userId, err)
		}

		if vsIPTransitDynamic.UserID == 0 {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找id:%d主账号信息不存在", userId)
		}

		flowRecordData := &model.FlowRecordData{}
		if vsIPTransitDynamic.FlowRecord != "" {
			json.Unmarshal([]byte(vsIPTransitDynamic.FlowRecord), flowRecordData)
		}

		recordsResult := []*model.VsIPFlowRecords{}

		if recordId == 0 {
			err = tx.
				Model(&model.VsIPFlowRecords{}).
				Select("id,flow,is_datacenter").
				Where(
					"user_id = ? and is_deal = 0 and  id >=  ?",
					userId, flowRecordData.LastId,
				).
				Order("id asc").
				Limit(limit).
				Find(&recordsResult).Error
			if err != nil {
				return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part1 err:%+v", err)
			}

		} else if recordId < flowRecordData.LastId {
			err = tx.
				Model(&model.VsIPFlowRecords{}).
				Select("id,flow,is_datacenter").
				Where(
					"user_id = ? and is_deal = 0 and  id  BETWEEN ? and   ?",
					userId, recordId, flowRecordData.LastId,
				).
				Order("id asc").
				Limit(limit).
				Find(&recordsResult).Error
			if err != nil {
				return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part2 err:%+v", err)
			}

		} else {
			err = tx.
				Model(&model.VsIPFlowRecords{}).
				Select("id,flow,is_datacenter").
				Where(
					"user_id = ? and is_deal = 0 and  id >=  ?",
					userId, recordId,
				).
				Order("id asc").
				Limit(limit).
				Find(&recordsResult).Error

			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part3 err:%+v", err)
		}

		recordsIds := []string{}

		for _, v := range recordsResult {
			if v.IsDatacenter {
				flowDatacenter += v.Flow
			} else {
				flow += v.Flow
			}
			recordsIds = append(recordsIds, fmt.Sprint(v.ID))
		}

		flow = int64(float64(flow) * flowIncRate)

		if flowDatacenter == 0 && flow == 0 && len(recordsResult) == 0 {
			return nil
		}

		flowRecordData.LastUseFlow = vsIPTransitDynamic.UseFlow
		flowRecordData.Flow = flow
		flowRecordData.LastResidualFlow = vsIPTransitDynamic.ResidualFlow

		flowRecordData.LastUseFlowDatacenter = vsIPTransitDynamic.UseFlowDatacenter
		flowRecordData.FlowDatacenter = flowDatacenter
		flowRecordData.LastResidualFlowDatacenter = vsIPTransitDynamic.ResidualFlowDatacenter

		if recordsResult[len(recordsResult)-1].ID >= flowRecordData.LastId {
			flowRecordData.RecordUnix = recordsResult[len(recordsResult)-1].Unix
			flowRecordData.LastId = recordsResult[len(recordsResult)-1].ID
		}

		flowRecordDataByte, _ := json.Marshal(flowRecordData)

		useFlow := vsIPTransitDynamic.UseFlow + uint64(flow)
		residualFlow := vsIPTransitDynamic.ResidualFlow - flow

		useFlowDatacenter := vsIPTransitDynamic.UseFlowDatacenter + uint64(flowDatacenter)
		residualFlowDatacenter := vsIPTransitDynamic.ResidualFlowDatacenter - flowDatacenter

		////修改总账号的流量
		dbRes := tx.Exec(
			`UPDATE`+vsIPTransitDynamic.TableName()+`
					SET
						use_flow =   ? ,
						residual_flow =   ? ,
						use_flow_datacenter = ?,
						residual_flow_datacenter = ?,
						flow_record = ?
							WHERE
								user_id = ? and
									use_flow = ? and
										residual_flow = ? and
											use_flow_datacenter = ? and
												residual_flow_datacenter = ?
			`,
			useFlow,
			residualFlow,
			useFlowDatacenter,
			residualFlowDatacenter,
			string(flowRecordDataByte),
			userId,
			vsIPTransitDynamic.UseFlow,
			vsIPTransitDynamic.ResidualFlow,
			vsIPTransitDynamic.UseFlowDatacenter,
			vsIPTransitDynamic.ResidualFlowDatacenter,
		)

		errTransaction := dbRes.Error

		if errTransaction != nil {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新主账号流量失败 err:%+v", errTransaction)
		}

		if dbRes.RowsAffected == 0 {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新主账号流量失败影响数据为0")
		}

		dbRes = tx.Exec(
			`UPDATE`+model.VsIPFlowRecords{}.TableName()+`
				SET
				is_deal = '1',
				deal_datetime = ?
					WHERE  id in  ?
			`,
			time.Now(),
			recordsIds,
		)
		errTransaction = dbRes.Error
		if errTransaction != nil {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新流量记录状态失败 err:%+v", errTransaction)
		}

		if dbRes.RowsAffected == 0 || int(dbRes.RowsAffected) != len(recordsIds) {
			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新流量记录状态失败数量为%d", dbRes.RowsAffected)
		}

		return nil
	})
	return
}
