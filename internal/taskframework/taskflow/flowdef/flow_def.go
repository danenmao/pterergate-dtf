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
