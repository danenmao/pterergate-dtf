package schedulingqueue

const (
	// 调度队列授予任务的基础时间片的大小, ms
	QueueBaseTimeSlice uint32 = 1000
	QueueTimeSliceStep uint32 = 1000
	RRQueueTimeSlice   uint32 = 5000

	// 优先级时间片bonus系数
	PriorityBonus_Low    = 1
	PriorityBonus_Medium = 8
	PriorityBonus_High   = 256

	// 任务调度静默期的上限
	QuietTaskMaxInterval = 600

	// 调度队列组中优先级队列的数量
	PrioirtyQueueCount uint32 = 2

	// 任务在优先级队列中分配的基础时间片的数量
	PriorityBaseQueueSliceCount uint32 = 40

	// 对优先级调度队列执行Priority Boost策略的间隔, 秒
	PriorityBoostInterval uint32 = 120
	PriorityBoostStep     uint32 = 60

	// 对低优先级调度队列执行Priority Boost策略的间隔, 秒
	RRPriorityBoostInterval uint32 = 240

	// 执行任务剩余时间加速策略的间隔, 秒
	RemainTaskAccelerationInteral uint32 = 300

	// 调度队列的Redis key前缀
	ScheduleQueueKeyPrefix = "Schedule"
)
