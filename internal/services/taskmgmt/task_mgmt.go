package taskmgmt

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/dbdef"
)

// 创建任务
func CreateTask(taskType uint32, param *taskmodel.TaskParam) (taskmodel.TaskIdType, error) {
	// 获取任务ID
	taskId, err := generateTaskId()
	if err != nil {
		glog.Warning("failed to create a task id: ", err)
		return 0, errordef.ErrOperationFailed
	}

	// 创建任务记录
	var taskRecord = dbdef.DBTaskRecord{}
	initTaskRecord(&taskRecord)
	taskRecord.Id = uint64(taskId)

	// 保存到taskInfo key中
	go saveInitTaskRecord(&taskRecord)

	// 向MySQL中添加任务记录
	err = addTaskRecord(taskId, param, &taskRecord)
	if err != nil {
		glog.Warning("failed to add task record: ", taskId, err)
		return 0, nil
	}

	// 启动创建协程
	go TaskCreationRoutine(taskId, taskType, param)

	glog.Info("succeeded to create a task: ", taskId)
	return taskId, nil
}

// 暂停任务
func PauseTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 恢复暂停中的任务
func ResumeTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 停止正在运行中的任务
func CancelTask(taskId taskmodel.TaskIdType) error {
	return nil
}

// 查询任务的运行状态
func GetTaskStatus(taskId taskmodel.TaskIdType, status *taskmodel.TaskStatusData) error {
	return nil
}
