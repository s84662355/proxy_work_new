package scheduled_tasks

import (
	"context"
	"sync"
	"time"

	"mproxy/constant"
	"mproxy/log"
	"mproxy/service"
	"mproxy/utils/taskConsumerManager"
)

var NewManager = sync.OnceValue(func() *manager {
	return &manager{
		tcm: taskConsumerManager.New(),
	}
})

type manager struct {
	tcm *taskConsumerManager.Manager
}

func (m *manager) Start() error {
	m.tcm.AddTask(1, func(ctx context.Context) {
		ticker := time.NewTicker(constant.DynamicAccountBatchUpdateDataLoopTime)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Info("[scheduled_tasks] 定时执行 BatchUpdateDynamicAccountData 上下文Done() 退出")
				return
			case t := <-ticker.C:
				log.Infof("[scheduled_tasks] 定时执行 BatchUpdateDynamicAccountData 开始 时间%+v", t)
				c, err := service.BatchUpdateDynamicAccountData(ctx)
				if err != nil {
					log.Errorf("[scheduled_tasks] 定时执行 BatchUpdateDynamicAccountData 时间%+v 数据数量%d 执行错误 err:%+v", t, c, err)
				}
				log.Infof("[scheduled_tasks] 定时执行 BatchUpdateDynamicAccountData 结束 时间%+v  数据数量%d", t, c)
			}
		}
	})

	return nil
}

func (m *manager) Stop() {
	m.tcm.Stop()
}
