package taskConsumerManager

import (
	"context"
	"sync"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:    ctx,
		cancel: cancel,
	}
}

type taskFunc func(context.Context)

func (m *Manager) AddTask(count int, fc taskFunc) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		tChan := make(chan struct{}, count)
		defer close(tChan)
		wg := sync.WaitGroup{}
		for {
			select {
			case <-m.ctx.Done():
				wg.Wait()
				return
			case tChan <- struct{}{}:
				wg.Add(1)
				go func() {
					defer func() {
						wg.Done()
						<-tChan
					}()
					fc(m.ctx)
				}()
			}
		}
	}()
}

func (m *Manager) Stop() {
	m.cancel()
	m.wg.Wait()
}
