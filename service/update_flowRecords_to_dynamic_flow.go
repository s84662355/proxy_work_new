package service

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// 	"gorm.io/gorm"
// 	"mproxy/constant"
// 	"mproxy/model"
// )

// func UpdateFlowRecordsToDynamicFlow(
// 	ctx context.Context,
// 	db *gorm.DB,
// 	userId int64,
// 	endTime int64,
// 	limit int,
// 	recordId uint64,
// ) (flow, flowDatacenter int64, err error) {
// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		vsIPTransitDynamic := &model.VsIPTransitDynamic{}
// 		err = tx.Where("user_id = ?  ", userId).First(vsIPTransitDynamic).Error
// 		if err != nil {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找id:%d主账号信息错误 err:%+v", userId, err)
// 		}

// 		if vsIPTransitDynamic.UserID == 0 {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找id:%d主账号信息不存在", userId)
// 		}

// 		flowRecordData := &model.FlowRecordData{}
// 		if vsIPTransitDynamic.FlowRecord != "" {
// 			json.Unmarshal([]byte(vsIPTransitDynamic.FlowRecord), flowRecordData)
// 		}

// 		recordsResult := []*model.VsIPFlowRecords{}

// 		if recordId == 0 {
// 			err = tx.
// 				Model(&model.VsIPFlowRecords{}).
// 				Select("id,flow").
// 				Where(
// 					"user_id = ? and is_deal = 0 and  id >=  ?",
// 					userId, flowRecordData.LastId,
// 				).
// 				Order("id asc").
// 				Limit(limit).
// 				Find(&recordsResult).Error
// 			if err != nil {
// 				return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part1 err:%+v", err)
// 			}

// 		} else if recordId < flowRecordData.LastId {
// 			err = tx.
// 				Model(&model.VsIPFlowRecords{}).
// 				Select("id,flow").
// 				Where(
// 					"user_id = ? and is_deal = 0 and  id  BETWEEN ? and   ?",
// 					userId, recordId, flowRecordData.LastId,
// 				).
// 				Order("id asc").
// 				Limit(limit).
// 				Find(&recordsResult).Error
// 			if err != nil {
// 				return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part2 err:%+v", err)
// 			}

// 		} else {
// 			err = tx.
// 				Model(&model.VsIPFlowRecords{}).
// 				Select("id,flow").
// 				Where(
// 					"user_id = ? and is_deal = 0 and  id >=  ?",
// 					userId, recordId,
// 				).
// 				Order("id asc").
// 				Limit(limit).
// 				Find(&recordsResult).Error

// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 查找流量记录失败 part3 err:%+v", err)
// 		}

// 		recordsIds := []string{}

// 		for _, v := range recordsResult {
// 			if v.IsDatacenter {
// 				flowDatacenter += v.Flow
// 			} else {
// 				flow += v.Flow
// 			}
// 			recordsIds = append(recordsIds, fmt.Sprint(v.ID))
// 		}

// 		flow = int64(float64(flow) * flowIncRate)

// 		if flowDatacenter == 0 && flow == 0 && len(recordsResult) == 0 {
// 			return nil
// 		}

// 		flowRecordData.LastUseFlow = vsIPTransitDynamic.UseFlow
// 		flowRecordData.Flow = flow
// 		flowRecordData.LastResidualFlow = vsIPTransitDynamic.ResidualFlow

// 		flowRecordData.LastUseFlowDatacenter = vsIPTransitDynamic.UseFlowDatacenter
// 		flowRecordData.FlowDatacenter = flowDatacenter
// 		flowRecordData.LastResidualFlowDatacenter = vsIPTransitDynamic.ResidualFlowDatacenter

// 		if recordsResult[len(recordsResult)-1].ID >= flowRecordData.LastId {
// 			flowRecordData.RecordUnix = recordsResult[len(recordsResult)-1].Unix
// 			flowRecordData.LastId = recordsResult[len(recordsResult)-1].ID
// 		}

// 		flowRecordDataByte, _ := json.Marshal(flowRecordData)

// 		useFlow := vsIPTransitDynamic.UseFlow + uint64(flow)
// 		residualFlow := vsIPTransitDynamic.ResidualFlow - flow

// 		useFlowDatacenter := vsIPTransitDynamic.UseFlowDatacenter + uint64(flowDatacenter)
// 		residualFlowDatacenter := vsIPTransitDynamic.ResidualFlowDatacenter - flowDatacenter

// 		log.Error("数据内容是---------------", map[string]interface{}{
// 			"residual_flow_datacenter": residualFlowDatacenter,
// 		})

// 		////修改总账号的流量
// 		dbRes := tx.Exec(
// 			`UPDATE`+vsIPTransitDynamic.TableName()+`
// 					SET
// 						use_flow =   ? ,
// 						residual_flow =   ? ,
// 						use_flow_datacenter = ?,
// 						residual_flow_datacenter = ?,
// 						flow_record = ?
// 							WHERE
// 								user_id = ? and
// 									use_flow = ? and
// 										residual_flow = ? and
// 											use_flow_datacenter = ? and
// 												residual_flow_datacenter = ?
// 			`,
// 			useFlow,
// 			residualFlow,
// 			useFlowDatacenter,
// 			residualFlowDatacenter,
// 			string(flowRecordDataByte),
// 			userId,
// 			vsIPTransitDynamic.UseFlow,
// 			vsIPTransitDynamic.ResidualFlow,
// 			vsIPTransitDynamic.UseFlowDatacenter,
// 			vsIPTransitDynamic.ResidualFlowDatacenter,
// 		)

// 		errTransaction := dbRes.Error

// 		if errTransaction != nil {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新主账号流量失败 err:%+v", errTransaction)
// 		}

// 		if dbRes.RowsAffected == 0 {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新主账号流量失败影响数据为0")
// 		}

// 		dbRes = tx.Exec(
// 			`UPDATE`+model.VsIPFlowRecords{}.TableName()+`
// 				SET
// 				is_deal = '1',
// 				deal_datetime = ?
// 					WHERE  id in  ?
// 			`,
// 			time.Now(),
// 			recordsIds,
// 		)
// 		errTransaction = dbRes.Error
// 		if errTransaction != nil {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新流量记录状态失败 err:%+v", errTransaction)
// 		}

// 		if dbRes.RowsAffected == 0 || int(dbRes.RowsAffected) != len(recordsIds) {
// 			return fmt.Errorf("UpdateFlowRecordsToDynamicFlow 更新流量记录状态失败数量为%d", dbRes.RowsAffected)
// 		}

// 		return nil
// 	})
// }
