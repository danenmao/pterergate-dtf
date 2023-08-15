package dtf

import (
	"pterergate-dtf/dtf/taskplugin"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/taskframework/taskloader"
)


////////////////////////////////////////////////////////////////////////
//
// 服务控制
//
////////////////////////////////////////////////////////////////////////

type ServiceRole uint32
const (
	ServiceRole_Manager ServiceRole = 1
	ServiceRole_Generator ServiceRole = 2
	ServiceRole_Scheduler ServiceRole = 3
	ServiceRole_Executor ServiceRole = 4
)


//
// 启动指定的服务
//
func StartService(role ServiceRole) error {
	return nil
}


//
// 通知停止服务
//
func NotifyStop() error {
	return nil
}


//
// 等待服务停止
//
func Join() error {
	return nil
}


////////////////////////////////////////////////////////////////////////
//
// 任务控制
//
////////////////////////////////////////////////////////////////////////

//
// 注册任务类型插件
//
func RegisterTaskType(register *taskplugin.TaskPluginRegister) error{
	return taskloader.RegisterTaskType(register)
}


//
// 创建任务
//
func CreateTask(taskType * uint32, param *taskmodel.TaskParam) (taskmodel.TaskIdType, error) {
	return 0, nil
}


//
// 暂停任务
//
func PauseTask(taskId taskmodel.TaskIdType) error {
	return nil
}


//
// 恢复暂停中的任务
//
func ResumeTask(taskId taskmodel.TaskIdType) error {
	return nil
}


//
// 停止正在运行中的任务
//
func CancelTask(taskId taskmodel.TaskIdType) error {
	return nil
}


//
// 查询任务的运行状态
//
func GetTaskStatus(taskId taskmodel.TaskIdType, status *taskmodel.TaskStatusData) error {
	return nil
}

