package taskplugin

import (
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

// 任务插件接口, 用于获取任务的执行信息
type ITaskPlugin interface {

	// 获取任务的配置信息
	GetTaskConf(taskConf *taskmodel.TaskConf) error

	// 获取任务的执行体
	GetTaskBody(taskBody *taskmodel.TaskBody) error
}

// 函数类型，用于返回指定任务类型所对应的任务插件对象
// 每个任务类型对应一个插件对象
type TaskPluginFactoryFn func(plugin *ITaskPlugin) error
