package taskloader

import (
	"sync"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
)

// 任务类型注册表
type PluginRegister struct {
	RegistrationTable map[uint32]*taskplugin.TaskPluginRegistration
	Lock              sync.Mutex
}

var gs_PluginRegister = PluginRegister{
	RegistrationTable: map[uint32]*taskplugin.TaskPluginRegistration{},
}

// 注册任务类型插件
func RegisterTaskType(register *taskplugin.TaskPluginRegistration) error {

	gs_PluginRegister.Lock.Lock()
	defer gs_PluginRegister.Lock.Unlock()

	_, ok := gs_PluginRegister.RegistrationTable[register.TaskType]
	if ok {
		glog.Info("found an existing task type plugin: ", register.TaskType)
		return nil
	}

	elem := *register
	gs_PluginRegister.RegistrationTable[register.TaskType] = &elem
	glog.Info("succeeded to register a task type: ", elem.TaskType)
	return nil
}
