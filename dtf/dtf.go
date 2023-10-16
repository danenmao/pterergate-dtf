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

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, opts ...ServiceOption) error {

	cfg := dtfdef.ServiceConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	return servicectrl.StartService(role, &cfg)
}

// 通知停止服务
func NotifyStop() error {
	return servicectrl.NotifyStop()
}

// 等待服务停止
func Join() error {
	return servicectrl.Join()
}

////////////////////////////////////////////////////////////////////////
//
// Task Type
//
////////////////////////////////////////////////////////////////////////

// 注册任务类型插件
func RegisterTaskType(register *taskplugin.TaskPluginRegister) error {
	return taskloader.RegisterTaskType(register)
}

////////////////////////////////////////////////////////////////////////
//
// Task Control
//
////////////////////////////////////////////////////////////////////////

// 创建任务
func CreateTask(taskType uint32, param *taskmodel.TaskParam) (taskmodel.TaskIdType, error) {
	return taskmgmt.CreateTask(taskType, param)
}

// 暂停任务
func PauseTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.PauseTask(taskId)
}

// 恢复暂停中的任务
func ResumeTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.ResumeTask(taskId)
}

// 停止正在运行中的任务
func CancelTask(taskId taskmodel.TaskIdType) error {
	return taskmgmt.CancelTask(taskId)
}

// 查询任务的运行状态
func GetTaskStatus(taskId taskmodel.TaskIdType, status *taskmodel.TaskStatusData) error {
	return taskmgmt.GetTaskStatus(taskId, status)
}
