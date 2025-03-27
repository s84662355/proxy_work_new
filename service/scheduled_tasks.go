//go:build scheduled_tasks
// +build scheduled_tasks

package service

import (
	"mproxy/service/scheduled_tasks"
)

func Start() error {
	return scheduled_tasks.NewManager().Start()
}

func Stop() {
	scheduled_tasks.NewManager().Stop()
}
