package taskmodel

import "time"

// PluginConf
// 记录任务插件的配置信息
type PluginConf struct {
	IterationMode   TaskInterationMode // 任务支持的迭代模式
	TaskTypeTimeout time.Duration      // 此任务类型的最大执行时间限制
}

// PluginBody
// 表示任务的执行体
type PluginBody struct {
	Generator         ITaskGenerator
	Executor          ITaskExecutor
	SchedulerCallback ITaskSchedulerCallback
	CollectorCallback ITaskCollectorCallback
}
