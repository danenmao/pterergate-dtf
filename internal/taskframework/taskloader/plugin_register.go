package taskloader

import (
	"sync"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
)

// 任务类型注册表
var gs_PluginRegisterTable = map[uint32]*taskplugin.TaskPluginRegister{}
var gs_PluginRegisterLock sync.Mutex

// 注册任务类型插件
func RegisterTaskType(register *taskplugin.TaskPluginRegister) error {

	gs_PluginRegisterLock.Lock()
	defer gs_PluginRegisterLock.Unlock()

	_, ok := gs_PluginRegisterTable[register.TaskType]
	if ok {
		glog.Info("found an existing task type plugin: ", register.TaskType)
		return nil
	}

	elem := *register
	gs_PluginRegisterTable[register.TaskType] = &elem
	glog.Info("succeeded to register a task type: ", elem.TaskType)
	return nil
}
