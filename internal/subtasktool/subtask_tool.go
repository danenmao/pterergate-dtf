package subtasktool

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/tasktool"
)

func SetSubtaskResult(
	subtaskId uint64,
	completeCode taskmodel.SubtaskResultType,
	scanResult interface{},
	ppipeline *redis.Pipeliner,
) error {

	pipeline := *ppipeline

	// get the task start time, calc the time cost of this subtask
	var startTime uint64 = 0
	var endTime uint64 = uint64(time.Now().Unix())
	var timeCost uint32 = 0
	var taskId taskmodel.TaskIdType = 0
	var appId uint32 = 0
	var taskType uint32 = 0

	err := tasktool.ReadSubtaskStartTime(subtaskId, &startTime, &taskId, &appId, &taskType)
	if err != nil {
		startTime = 0
		timeCost = 1
	} else {
		timeCost = uint32(endTime - startTime)
	}

	// determine the subtask status value
	subtaskStatus := taskmodel.SubtaskStatus_Finished
	if completeCode == taskmodel.SubtaskResult_Timeout {
		subtaskStatus = taskmodel.SubtaskStatus_Timeout
	}

	// save subtask result to subtask_info key
	values := map[string]interface{}{
		config.SubtaskInfo_SubtaskResult: scanResult,
		config.SubtaskInfo_EndTimeField:  endTime,
		config.SubtaskInfo_TimeCostField: timeCost,
		config.SubtaskInfo_Complete_code: completeCode,
		config.SubtaskInfo_StatusField:   subtaskStatus,
	}

	pipeline.HMSet(
		context.Background(),
		tasktool.GetSubtaskKey(subtaskId),
		values,
	)

	// inc subtask completion reason counter
	if taskId != 0 {
		subtaskFieldMap := map[taskmodel.SubtaskResultType]string{
			taskmodel.SubtaskResult_Success: config.TaskInfo_CompletedSubtaskCountField,
			taskmodel.SubtaskResult_Failure: config.TaskInfo_CompletedSubtaskCountField,
			taskmodel.SubtaskResult_Timeout: config.TaskInfo_TimeoutSubtaskCountField,
		}

		fieldName, ok := subtaskFieldMap[completeCode]
		if !ok {
			glog.Error("invalid complete code: ", subtaskId, taskId, completeCode)
		} else {
			pipeline.HIncrBy(context.Background(), tasktool.GetTaskInfoKey(taskId), fieldName, 1)
		}
	}

	return nil
}

func ReadSubtaskStatus(subtaskId taskmodel.SubtaskIdType, statusRet *uint32) error {

	cmd := redistool.DefaultRedis().HGet(context.Background(), tasktool.GetSubtaskKey(uint64(subtaskId)),
		config.SubtaskInfo_StatusField)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get status of subtask: ", subtaskId, err)
		return err
	}

	val := cmd.Val()

	status, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		glog.Warning("failed to parse status field: ", subtaskId, val, err)
		return err
	}

	*statusRet = uint32(status)
	return nil
}

func GetSubtaskTaskType(subtaskId uint64, retTaskType *uint32) error {
	return GetSubtaskUint(subtaskId, config.SubtaskInfo_TaskTypeField, retTaskType)
}

// 取subtask info key中记录的uint数据
func GetSubtaskUint(subtaskId uint64, field string, retVal *uint32) error {

	strCmd := redistool.DefaultRedis().HGet(context.Background(), tasktool.GetSubtaskKey(subtaskId), field)
	err := strCmd.Err()
	if err != nil {
		glog.Warning("failed to get field of subtask: ", subtaskId, ",", field, ",", err)
		return err
	}

	val, err := strCmd.Uint64()
	if err != nil {
		glog.Warning("failed to convert field of subtask: ", subtaskId, ",", field, ",", strCmd.Val(), err)
		return err
	}

	*retVal = uint32(val)
	return nil
}

func IsSubtaskRunning(subtaskId taskmodel.SubtaskIdType) bool {

	var status uint32 = 0
	err := ReadSubtaskStatus(subtaskId, &status)
	if err != nil {
		glog.Warning("failed to read subtask status: ", subtaskId, err)
		return false
	}

	return status == taskmodel.SubtaskStatus_Running
}
