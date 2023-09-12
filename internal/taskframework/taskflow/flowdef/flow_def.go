package flowdef

// 任务的创建参数
type TaskCreateParam struct {
	ResourceGroupName string `json:"resource_group"`
	TaskType          uint32 `json:"task_type"`
	Priority          uint32 `json:"priority"`
	Timeout           uint32 `json:"timeout"`
	TypeParam         string `json:"type_param"`
}

// 保存任务创建
const TaskCreateParamPrefix = "task.create.param."

// 任务调度数据
type TaskScheduleData struct {
	ResourceGroupName   string `json:"resource_group"`         // 任务所属的资源组名
	CurrentQueue        uint32 `json:"current_queue"`          // 任务所属的调度队列索引, 0开始编号
	CurrentQueueKeyName string `json:"current_queue_key_name"` // 任务所属的调度队列名
	InitiallQueueSlice  uint32 `json:"initial_queue_slice"`    // 任务在当前调度队列中的初始时间片数量
	QueueSlice          uint32 `json:"queue_slice"`            // 任务在当前调度队列中的时间片数量
	QuietStartTime      uint64 `json:"quiet_start_time"`       // 任务静默的起始时间
}

// 保存任务调度数据
const RedisTaskScheduleDataPrefix = "task.schedule.data."

// 执行器服务名
const SubtaskExecutorName = ""
