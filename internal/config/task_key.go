package config

import "time"

// 默认的key有效期
const DefaultExpiredValue = time.Hour * 24 * 7

// 任务键定义
const (
	// 任务ID计数器
	TaskIdKey = "dtf.task.id.counter"

	// 子任务ID计数器
	SubtaskIdKey = "dtf.subtask.id.counter"
)

const (

	// task_zset
	TaskZset = "dtf.task.list"

	Stage_CreatingTask = "stage_creating_task"

	// 创建中的任务的集合, creating_task_zset
	CreatingTaskZset = "dtf.creating.task.list"

	// 生成中的任务的有序集合, task_generation_zset
	GeneratingTaskZset = "dtf.generating.task.list"

	// 待生成的任务集合, task_to_generate_zset
	ToGenerateTaskZset = "dtf.to.generate.task.list"

	// 执行中的任务集合, task_schedule_zset
	RunningTaskZset = "dtf.running.task.list"

	// 取消中的任务集合, cancelling_task_zset
	CancellingTaskList = "dtf.cancel.task.list"

	// 暂停中的任务集合, pausing_task_zset
	PausingTaskList = "dtf.pause.task.list"

	// 已完成待处理的任务集体, completed_task_zset
	CompletedTaskList = "dtf.completed.task.list"
)

const (
	// 任务信息key, task_info.$taskid
	TaskInfoKeyPrefix                   = "dtf.task.info."
	TaskInfo_StageField                 = "stage"
	TaskInfo_StepField                  = "step"
	TaskInfo_UID                        = "uid"
	TaskInfo_CreateTimeField            = "create_time"
	TaskInfo_TotalSubtaskCountField     = "total_subtask_count"
	TaskInfo_CompletedSubtaskCountField = "completed_subtask_count"
	TaskInfo_TimeoutSubtaskCountField   = "timeout_subtask_count"
	TaskInfo_CancelledSubtaskCountField = "cancelled_subtask_count"
	TaskInfo_GenerationCompletedField   = "generation_completed"
	TaskInfo_ResourceCostField          = "resource_cost"
	TaskInfo_TaskTypeField              = "task_type"
	TaskInfo_Progess                    = "progress"
	TaskInfo_TypeParam                  = "type_param"
	TaskInfo_InitTaskRecord             = "init_task_record"
	TaskInfo_CheckUIDMapField           = "check_uid_map"
	TaskInfo_StatusField                = "status" // 任务的运行状态: 1:运行中; 2:已完成; 3:已取消;

	// 每个任务的锁
	TaskInfoLockPrefix = "dtf.task.lock."
)

const (
	// next_check_time的定义
	NextCheckTimeField = "next_check_time"

	// 任务生成阶段的进度key, task_generation.$taskid.progress
	TaskGenerationKeyPrefix              = "dtf.task.generation.progress."
	TaskGenerationKey_StepField          = "step"
	TaskGenerationKey_NextCheckTimeField = "next_check_time"
)

const (
	// 每个客户的锁, user_concurrency_lock.$uid
	UserNodeLockPrefix = "dtf.user.node.lock."

	// 任务下的节点的有序集合, node_list_$taskId
	TaskNodeListPrefix = "dtf.task.node.list."

	// 任务下当前要处理的节点的索引，node_list_pointer_$taskid
	TaskNextNodePointerPrefix = "dtf.task.next.node.pointer."

	// 任务已下发的子任务的列表，pushed_subtask_list_$taskid
	TaskPushedSubtaskList = "dtf.pushed.subtask.list."

	// 任务下的子任务集合, subtask_list.$taskid, set
	TaskToSubtaskSetPrefix = "dtf.task.subtask.list."
)

const (
	// 待调度的子任务队列
	ToScheduleSubtaskZset = "dtf.subtask.to.schedule.list"

	// 执行中的子任务队列
	RunningSubtaskZset = "dtf.running.subtask.list"

	// 已完成的子任务的集合, zset, 按照完成时间排序
	CompletedSubtaskList = "dtf.completed.subtask.list"

	// 正在重试的子任务的集合
	RetryTimeoutSubtaskList = "dtf.retry.timeout.subtask.list"

	// 正准备重试的子任务的有序集合
	ToRetryTimeoutSubtaskZset = "dtf.to.retry.timeout.subtask.list"
)

const (
	// 子任务的信息
	SubtaskInfoPrefix             = "dtf.subtask.info."
	SubtaskInfo_UID               = "uid"
	SubtaskInfo_TaskIdField       = "task_id"          // 子任务所属的任务ID
	SubtaskInfo_TaskTypeField     = "task_type"        // 子任务的任务类型
	SubtaskInfo_TimeoutField      = "timeout"          // 子任务的超时时间
	SubtaskInfo_TimeoutCountField = "timeout_count"    // 子任务超时的次数
	SubtaskInfo_PriorityField     = "subtask_priority" // 子任务的优先级
	SubtaskInfo_StartTimeField    = "start_time"       // 子任务执行的开始时间
	SubtaskInfo_EndTimeField      = "end_time"         // 子任务执行的结束时间
	SubtaskInfo_TimeCostField     = "time_cost"        // 子任务的执行耗时
	SubtaskInfo_Complete_code     = "complete_code"    // 子任务的完成码
	SubtaskInfo_SubtaskResult     = "subtask_result"   // 执行结果
	SubtaskInfo_Param             = "param"            // 子任务的执行参数
	SubtaskInfo_StatusField       = "status"           // 子任务的运行状态

)

const (
	// 任务与其包含的对象对应关系
	TaskObjectSetPrefix = "dtf.task.object.set."

	// 任务中的对象与多个子任务的对应关系
	TaskObjectSubtaskSetPrefix = "dtf.object.subtask.set."
)
