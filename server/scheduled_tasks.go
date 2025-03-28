//go:build scheduled_tasks
// +build scheduled_tasks

package server

import (
	"mproxy/server/scheduled_tasks"
)

func Start() error {
	return scheduled_tasks.NewManager().Start()
}

func Stop() {
	scheduled_tasks.NewManager().Stop()
}
