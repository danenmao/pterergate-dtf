package dtf

import (
	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
	"github.com/danenmao/pterergate-dtf/internal/servicectrl"
	"github.com/danenmao/pterergate-dtf/internal/services/taskmgmt"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskloader"
)

////////////////////////////////////////////////////////////////////////
//
// Service Control
//
////////////////////////////////////////////////////////////////////////

// start the specified service
func StartService(role dtfdef.ServiceRole, opts ...ServiceOption) error {
	cfg := dtfdef.ServiceConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	return servicectrl.StartService(role, &cfg)
}

// notify to stop the service
func NotifyStop() error {
	return servicectrl.NotifyStop()
}

// wait for the service to stop
func Join() error {
	return servicectrl.Join()
}

////////////////////////////////////////////////////////////////////////
//
// Task Type
//
////////////////////////////////////////////////////////////////////////

// register a task type plugin
func RegisterTaskType(register *taskplugin.TaskPluginRegistration) error {
	return taskloader.RegisterTaskType(register)
}

////////////////////////////////////////////////////////////////////////
//
// Task Control
//
////////////////////////////////////////////////////////////////////////

// create a task
func CreateTask(taskType uint32, param *taskmodel.TaskParam) (taskmodel.TaskIdType, error) {
	return taskmgmt.CreateTask(taskType, param)
}

// pause a running task
func PauseTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.PauseTask(taskId)
}

// resume a paused task
func ResumeTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.ResumeTask(taskId)
}

// cancel a running task
func CancelTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.CancelTask(taskId)
}

// retrieve the task status
func GetTaskStatus(taskId taskmodel.TaskIdType, status *taskmodel.TaskStatusData) error {
	return taskmgmt.GetTaskStatus(taskId, status)
}
