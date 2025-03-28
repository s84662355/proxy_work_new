//go:build scheduled_tasks
// +build scheduled_tasks

package config

const ServiceName = "scheduled_tasks"

type serviceConf struct{
	TaskCount  int `json:"task_count" `
}
