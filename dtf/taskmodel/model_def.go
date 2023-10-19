package taskmodel

import (
	"time"
)

// 迭代模式
type TaskInterationMode int

const (
	IterationMode_No           TaskInterationMode = 1 // 不支持迭代
	IterationMode_UseCollector TaskInterationMode = 2 // 支持子任务迭代
)

// 任务状态的类型
type TaskStatusType uint32

const (
	TaskStatus_Created     TaskStatusType = 1 // 已创建
	TaskStatus_Running     TaskStatusType = 2 // 运行中
	TaskStatus_Paused      TaskStatusType = 3 // 暂停中
	TaskStatus_Cacelled    TaskStatusType = 4 // 已取消
	TaskStatus_Completed   TaskStatusType = 5 // 已完成
	TaskStatus_Exceptional TaskStatusType = 6 // 异常
)

// 任务创建者结构
type TaskCreator struct {
	UID  uint64 `json:"uid"`  // 用户id, 由调用者自行定义
	Name string `json:"name"` // 用户名, 同调用者自行定义
}

// 任务的执行状态数据
type TaskStatusData struct {
	TaskId        TaskIdType     `json:"task_id"`        // 任务ID
	TaskType      uint32         `json:"task_type"`      // 任务类型
	TaskStatus    TaskStatusType `json:"task_status"`    // 任务的状态
	TaskProgress  float32        `json:"task_progress"`  // 任务执行的进度
	Priority      uint32         `json:"priority"`       // 任务优先级
	SubtaskCount  uint32         `json:"subtask_count"`  // 任务包含的子任务数量
	StartTime     time.Time      `json:"start_time"`     // 任务的开始时间
	ResourceGroup string         `json:"resource_group"` // 任务所属的资源组名
	TaskName      string         `json:"task_name"`      // 任务名
}

// 任务的创建参数
type TaskParam struct {
	Creator       TaskCreator   `json:"creator"`        // 任务创建者
	ResourceGroup string        `json:"resource_group"` // 任务所属的资源组名称
	Priority      uint32        `json:"priority"`       // 任务的基础优先级
	TaskName      string        `json:"task_name"`      // 任务名
	Description   string        `json:"description"`    // 任务描述
	TaskType      uint32        `json:"task_type"`      // 任务类型
	Timeout       time.Duration `json:"timeout"`        // 任务的超时值
	TypeParam     string        `json:"type_param"`     // 任务的自定义参数
}

// 任务的执行结果
type TaskResult struct {
	TaskId     TaskIdType     `json:"task_id"`     // 任务ID
	Result     TaskStatusType `json:"result"`      // 任务的结果
	ResultCode uint32         `json:"result_code"` // 任务的结果码
	Reason     string         `json:"reason"`      // 结果描述
	ResultData string         `json:"result_data"` // 与任务类型相关的结果数据
}
