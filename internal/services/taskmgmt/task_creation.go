package taskmgmt

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/basedef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/dbdef"
	"github.com/danenmao/pterergate-dtf/internal/idtool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/flowdef"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

// 任务创建协程
func TaskCreationRoutine(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	taskParam *taskmodel.TaskParam,
) {
	glog.Info("begin to create a task: ", taskId)

	// 将 $taskid 添加到创建中的任务列表
	err := tasktool.AddTaskToCreatingQueue(taskId)
	if err != nil {
		glog.Warning("failed to add task to creating queue, return: ", err)
		return
	}

	// 为任务添加task info key, 为 hash key.
	err = tasktool.CreateTaskInfoKey(taskId, taskParam)
	if err != nil {
		glog.Warning("failed to add task to creating task list: ", taskId, err.Error())
		return
	}

	// 将任务添加到已存在任务列表中, 表示任务已经存在
	err = tasktool.AddTaskToExistingTaskList(taskId, taskParam.Timeout)
	if err != nil {
		glog.Warning("failed to add task to existing list, return: ", err)
		return
	}

	// 在Redis key中保存任务创建中指定的TypeParam
	tasktool.SetTaskRawTypeParam(taskId, taskParam.TypeParam)

	// 记录任务的创建参数
	createParam := flowdef.TaskCreateParam{
		ResourceGroupName: taskParam.ResourceGroup,
		TaskType:          taskParam.TaskType,
		Priority:          taskParam.Priority,
		Timeout:           uint32(taskParam.Timeout / time.Second),
		TypeParam:         taskParam.TypeParam,
	}
	tasktool.SaveTaskCreateParam(taskmodel.TaskIdType(taskId), &createParam)

	// 结束创建过程
	finishInitialization(taskId, taskType)
	glog.Info("succeeded to create a task, task creation routine exited: ", taskId)
}

// 生成任务ID
func generateTaskId() (taskmodel.TaskIdType, error) {

	id, err := idtool.GetAvailableId(config.TaskIdKey)
	if err != nil {
		glog.Warning("failed to get task id: ", err)
		return 0, errors.New("failed to get task id")
	}

	return taskmodel.TaskIdType(id), nil
}

// 初始化任务结构
func initTaskRecord(taskRecord *dbdef.TaskRecord) {

	taskRecord.Id = 0
	taskRecord.TaskStatus = uint8(taskmodel.TaskStatus_Running)
	taskRecord.UID = 0
	taskRecord.Creator = ""
	taskRecord.Name = ""
	taskRecord.Description = ""
	taskRecord.TaskType = 0

	taskRecord.StartTime = time.Now().Format(basedef.GoTimeFormatStr)
	taskRecord.FinishTime = dbdef.DBNullTimeStr
	taskRecord.TimeCost = 0

	// 设置NextCheckTime字段, 用于超时检查
	taskRecord.NextCheckTime = time.Now().Add(time.Second *
		time.Duration(config.EnvTaskCreationNextCheck)).Local().Format(dbdef.GoTimeFormatStr)
}

// 添加任务记录
func addTaskRecord(
	taskId taskmodel.TaskIdType,
	param *taskmodel.TaskParam,
	taskRecord *dbdef.TaskRecord,
) error {

	taskRecord.UID = param.Creator.UID
	taskRecord.Creator = param.Creator.Name
	taskRecord.TaskType = param.TaskType
	taskRecord.Name = param.TaskName
	taskRecord.Description = param.Description

	// 往数据库中添加任务记录
	err := tasktool.AddTaskRecord(taskRecord)
	if err != nil {
		glog.Warning("failed to add task record: ", err.Error())
		return err
	}

	return nil
}

// 保存任务的初始化结构
func saveInitTaskRecord(taskRecord *dbdef.TaskRecord) {

	if taskRecord == nil {
		panic("invalid task record pointer")
	}

	if taskRecord.Id == 0 {
		panic("invalid task id")
	}

	data, err := json.Marshal(taskRecord)
	if err != nil {
		glog.Warning("failed to marshal task record: ", taskRecord.Id, err)
		return
	}

	taskInfoKey := tasktool.GetTaskInfoKey(taskmodel.TaskIdType(taskRecord.Id))
	cmd := redistool.DefaultRedis().HSet(context.Background(), taskInfoKey,
		config.TaskInfo_InitTaskRecord, string(data))
	if cmd.Err() != nil {
		glog.Warning("failed to set init task record for task: ", taskRecord.Id, cmd.Err())
		return
	}

	glog.Info("succeeded to save init task record of task: ", taskRecord.Id)
}

// 完成任务初始化
func finishInitialization(
	taskId taskmodel.TaskIdType,
	taskType uint32,
) {
	glog.Info("succeeded to finish initialization of task: ", taskId)
}
