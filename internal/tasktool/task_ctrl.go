package tasktool

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/basedef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/dbdef"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

// 完成任务
func CompleteTask(taskId taskmodel.TaskIdType, taskRecord *dbdef.DBTaskRecord) error {

	// 从task info key中读取信息
	taskKey := GetTaskInfoKey(taskId)
	mapCmd := redistool.DefaultRedis().HGetAll(context.Background(), taskKey)
	infos, err := mapCmd.Result()
	if err != nil {
		glog.Warning("failed to get task info: ", taskId, err.Error())
		return err
	}

	// 取task type
	taskTypeStr, ok := infos[config.TaskInfo_TaskTypeField]
	if !ok {
		return errors.New("no task type field in task info key")
	}

	taskType, err := strconv.Atoi(taskTypeStr)
	if err != nil {
		glog.Warning("invalid task type: ", taskTypeStr, err.Error())
		return errors.New("invalid task type")
	}

	// create time为时间戳
	createTimeStr, ok := infos[config.TaskInfo_CreateTimeField]
	if !ok {
		glog.Warning("no create time field in task info key: ", taskId)
		return errors.New("no create time field in task info key")
	}

	createTime, err := strconv.ParseUint(createTimeStr, 10, 64)
	if err != nil {
		glog.Warning("failed to parse create time: ", taskId, createTimeStr, err)
		return errors.New("failed to parse create time")
	}

	timeCost := uint64(time.Now().Unix()) - createTime

	// uid
	udStr, ok := infos[config.TaskInfo_UID]
	if !ok {
		glog.Warning("no uid field in task info key: ", taskId)
		udStr = "0"
	}

	uid, err := strconv.ParseUint(udStr, 10, 32)
	if err != nil {
		glog.Warning("failed to parse appid: ", udStr, err)
		uid = 0
	}

	// 更新task info key, 写入完成状态
	SetTaskStatus(taskId, taskmodel.TaskStatus_Completed)

	// 更新 task 表的内容
	taskRecord.Id = uint64(taskId)
	taskRecord.TaskStatus = uint8(taskmodel.TaskStatus_Completed)
	taskRecord.FinishTime = time.Now().Format(basedef.GoTimeFormatStr)
	taskRecord.TimeCost = uint32(timeCost)
	taskRecord.TaskType = uint32(taskType)
	taskRecord.UID = uid

	err = WriteCompleteInfoToTaskDB(taskRecord)
	if err != nil {
		glog.Warning("failed to update task table: ", taskId, err.Error())
		return err
	}

	glog.Info("succeeded to complete task: ", taskId)
	return nil
}

// 设置任务的运行状态
func SetTaskStatus(taskId taskmodel.TaskIdType, status taskmodel.TaskStatusType) error {

	cmd := redistool.DefaultRedis().HSet(context.Background(), GetTaskInfoKey(taskId), config.TaskInfo_StatusField, status)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to set status of task: ", taskId, err)
		return err
	}

	return nil
}

// 更新任务表记录，标记任务已完成
func WriteCompleteInfoToTaskDB(taskRecord *dbdef.DBTaskRecord) error {
	result, err := mysqltool.DefaultMySQL().NamedExec(
		dbdef.SQL_TaskTable_CompleteTask,
		taskRecord,
	)

	if err != nil {
		glog.Warning("failed to update task result: ", taskRecord.Id, err.Error())
		return err
	}

	lines, _ := result.RowsAffected()
	glog.Info("succeeded to write complete info to task table: ", taskRecord.Id, lines)
	return nil
}
