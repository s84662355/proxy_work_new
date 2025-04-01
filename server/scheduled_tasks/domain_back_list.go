package scheduled_tasks

import (
	"context"
	"time"

	"go.uber.org/zap"
	"mproxy/common"
	"mproxy/config"
	"mproxy/dao"
	"mproxy/log"
	"mproxy/service"
)

// 推送域名黑名单到网关
func (m *manager) pushBlacklistToProxyGateWay(ctx context.Context) {
	loopTime := 10 * 60 * time.Second
	ticker := time.NewTicker(10 * 60 * time.Second)
	defer ticker.Stop()
	for {
		ticker.Reset(loopTime)
		select {
		case <-ctx.Done():
			log.Info("[定时任务scheduled_tasks]推送域名黑名单到网关任务上下文结束退出")
			return
		case <-ticker.C:
			log.Info("[定时任务scheduled_tasks]定时执行推送域名黑名单到网关开始")

			bl, err := dao.GetDomainBackListOnlyDomain(
				ctx,
				common.GetMysqlDB(),
			)
			if err != nil {
				log.Error("[定时任务scheduled_tasks]定时执行推送域名黑名单到网关获取域名黑名单错误", zap.Any("error", err))
			}

			if err := service.PushDomainBackList(
				common.GetRabbitMqProductPool(),
				config.GetConf().Rabbitmq.BlacklistExchange,
				bl,
			); err != nil {
				log.Error("[定时任务scheduled_tasks]定时执行推送域名黑名单到网关rabbitmq推送", zap.Any("error", err))
			}

		}
	}
}
