package taskloader

import (
	"errors"
	"sync"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskplugin"
)

// 插件对象表, 及锁
var gs_TaskPluginTable = map[uint32]taskplugin.ITaskPlugin{}
var gs_TaskPluginLock sync.Mutex

// 查找指定类型任务的插件对象
func LookupTaskPlugin(taskType uint32, plugin *taskplugin.ITaskPlugin) error {

	gs_TaskPluginLock.Lock()
	defer gs_TaskPluginLock.Unlock()

	// check if task plugin exists
	elem, ok := gs_TaskPluginTable[taskType]
	if ok {
		glog.Info("found task type plugin: ", taskType)
		*plugin = elem
		return nil
	}

	// create a task plugin
	err := NewTaskPlugin(taskType, plugin)
	if err != nil {
		glog.Error("failed to load task plugin: ", taskType, ",", err)
		return err
	}

	gs_TaskPluginTable[taskType] = *plugin
	glog.Info("succeeded to save a task plugin: ", taskType)
	return nil
}

// 加载任务插件, 返回任务插件指针
func NewTaskPlugin(taskType uint32, plugin *taskplugin.ITaskPlugin) error {

	// 查找传入任务类型对应的加载器
	register, ok := gs_PluginRegisterTable[taskType]
	if !ok {
		glog.Error("unknown task type: ", taskType)
		return errors.New("unknown task type")
	}

	if register.TaskType != taskType {
		glog.Error("unmatched task type: ", taskType)
		return errors.New("unmatched task type")
	}

	// 调用函数来创建任务插件
	err := register.NewPluginFn(plugin)
	if err != nil {
		glog.Error("failed to new a plugin: ", taskType, err)
		return err
	}

	glog.Info("succeeded to new a plugin: ", taskType, ", ", register.Description)
	return nil
}
