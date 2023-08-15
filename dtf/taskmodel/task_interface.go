package taskmodel

// 任务的生成接口
// 接口将任务分解为可以并行执行的子任务, 返回给任务框架
type ITaskGenerator interface {

	// 任务开始时，通知接口开始执行生成操作
	Begin(taskId TaskIdType, taskType uint32, taskData *TaskParam) error

	// 任务结束时，通知接口进行清理
	End(taskId TaskIdType) error

	// 通知接口取消生成操作
	Cancel(taskId TaskIdType) error

	// 通知接口暂停生成操作
	Pause(taskId TaskIdType) error

	// 通知接口恢复生成操作
	Resume(taskId TaskIdType) error

	// 查询任务的进度
	QueryProgress(taskId TaskIdType) (float32, error)

	// 获取下一个子任务
	// 若无子任务,返回errordef.ErrNotFound; 否则返回nil;
	// 若生成完成，设置finished
	GetSubtask(taskId TaskIdType, subtaskData *SubtaskData, finished *bool) error

	// 当每个子任务生成完成后，通知接口
	AfterGeneration(subtaskId SubtaskIdType, subtask *SubtaskData) error
}

// 任务的调度接口
// 响应一些调度操作
type ITaskScheduler interface {

	// 在子任务被调度之前调用，可通过返回的bool来控制当前子任务是否被调度
	BeforeDispatch(subtaskId SubtaskIdType, subtaskData *SubtaskData) (bool, error)

	// 子任务被调度进入执行后执行
	AfterDispatch(subtaskId SubtaskIdType) error
}

// 任务的执行接口
// 实现任务的主要操作
type ITaskExecutor interface {

	// 实现子任务的计算操作
	DoCalculation(subtaskData *SubtaskData, result *SubtaskResult) error

	// 通知接口退出
	Cancel() error
}

// 任务的迭代接口
type ITaskIterator interface {

	// 每次返回一次结果, 调用一次方法
	OnResult(subtaskId SubtaskIdType, subtaskResult *SubtaskResult) error

	// 一个子任务执行完成后执行
	AfterCalculation(subtaskId SubtaskIdType, subtaskResult *SubtaskResult) (int, error)

	// 整个任务完成时执行
	AfterTaskCompleted(taskId TaskIdType) (int, error)
}

// IDataTranlator
type IDataTranlator interface {
}