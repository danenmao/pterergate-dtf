package taskmodel

// TaskConf
// 记录任务插件的配置信息
type TaskConf struct {
	IterationMode    TaskInterationMode // 任务支持的迭代模式
	TypeTimeoutInSec uint32             // 此任务类型的最大执行时间限制
}

// TaskBody
// 表示任务的执行体
type TaskBody struct {
	Generator ITaskGenerator
	Scheduler ITaskScheduler
	Iterator  ITaskIterator
	Executor  ITaskExecutor
}
