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
	return nil
	// m.tcm.AddTask()
}

func (m *manager) Stop() {
	m.tcm.Stop()
}
