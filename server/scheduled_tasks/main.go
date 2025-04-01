package scheduled_tasks

import (
	"sync"

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
	m.tcm.AddTask(1, m.batchUpdateDynamicAccountDataCache)
	m.tcm.AddTask(1, m.batchUpdateDynamicDataCache)
	m.tcm.AddTask(2, m.updateDynamicAccountFlowFromRedisToDB)
	m.tcm.AddTask(1, m.checkFlowRecordsToRedisSortedSet)
	m.tcm.AddTask(2, m.updateFlowRecordsFromRedisSortedSet)
	m.tcm.AddTask(1, m.updateVsIpFlowRecordsPartition)

	return nil
}

func (m *manager) Stop() {
	m.tcm.Stop()
}
