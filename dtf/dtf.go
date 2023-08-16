package dtf

import (
	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/dtf/taskplugin"
	"pterergate-dtf/internal/servicectl"
	"pterergate-dtf/internal/taskframework/taskloader"
)

////////////////////////////////////////////////////////////////////////
//
// 服务控制
//
////////////////////////////////////////////////////////////////////////

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, opts ...ServiceOption) error {

	cfg := dtfdef.ServiceConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	return servicectl.StartService(role, &cfg)
}

// 通知停止服务
func NotifyStop() error {
	return servicectl.NotifyStop()
}

// 等待服务停止
func Join() error {
	return servicectl.Join()
}

////////////////////////////////////////////////////////////////////////
//
// 任务控制
//
////////////////////////////////////////////////////////////////////////

// 注册任务类型插件
func RegisterTaskType(register *taskplugin.TaskPluginRegister) error {
	return taskloader.RegisterTaskType(register)
}

// 创建任务
func CreateTask(taskType *uint32, param *taskmodel.TaskParam) (taskmodel.TaskIdType, error) {
	return 0, nil
}

// 暂停任务
func PauseTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 恢复暂停中的任务
func ResumeTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 停止正在运行中的任务
func CancelTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 查询任务的运行状态
func GetTaskStatus(taskId taskmodel.TaskIdType, status *taskmodel.TaskStatusData) error {
	return nil
}
