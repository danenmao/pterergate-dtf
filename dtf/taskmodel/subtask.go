package taskmodel

import (
	"time"
)

// 子任务结果类型
type SubtaskResultType uint32

const (
	SubtaskResult_Success SubtaskResultType = 1 // 子任务执行成功
	SubtaskResult_Failure SubtaskResultType = 2 // 子任务执行失败
	SubtaskResult_Timeout SubtaskResultType = 3 // 子任务超时
)

// 子任务的运行状态
const (
	SubtaskStatus_Running   = 1
	SubtaskStatus_Finished  = 2
	SubtaskStatus_Cancelled = 3
	SubtaskStatus_Timeout   = 4
)

// 子任务的参数体
type SubtaskBody struct {
	SubtaskId    SubtaskIdType `json:"subtask_id"`    // 子任务ID
	TaskId       TaskIdType    `json:"task_id"`       // 所属的任务ID
	TaskType     uint32        `json:"task_type"`     // 任务类型
	Timeout      uint32        `json:"timeout"`       // 子任务的超时值, 秒
	TypeParam    string        `json:"type_param"`    // 子任务与类型相关的参数数, json
	CreatedAt    time.Time     `json:"created_at"`    // 子任务创建的时间
	TerminatedAt time.Time     `json:"terminated_at"` // 子任务结束的时间
}

// 子任务执行的结果
type SubtaskResult struct {
	SubtaskId  SubtaskIdType     `json:"subtask_id"`  // 子任务ID
	TaskId     TaskIdType        `json:"task_id"`     // 所属的任务ID
	Result     SubtaskResultType `json:"result"`      // 子任务的结果
	ResultCode uint32            `json:"result_code"` // 子任务的结果码
	ResultMsg  string            `json:"result_msg"`  // 原因
	ResultBody string            `json:"result_body"` // 子任务与类型相关的结果数据
}
