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

	// 取子任务的开始时间
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

	// 将子任务执行结果保存到subtask_info.$subtaskid中
	values := map[string]interface{}{
		config.SubtaskInfo_SubtaskResult: scanResult,
		config.SubtaskInfo_EndTimeField:  endTime,
		config.SubtaskInfo_TimeCostField: timeCost,
		config.SubtaskInfo_Complete_code: completeCode,
		config.SubtaskInfo_StatusField:   taskmodel.SubtaskStatus_Finished,
	}

	pipeline.HMSet(
		context.Background(),
		tasktool.GetSubtaskKey(subtaskId),
		values,
	)

	// inc subtask completion reason counter
	if taskId != 0 {
		subtaskFieldMap := map[taskmodel.SubtaskResultType]string{
			taskmodel.SubtaskStatus_Finished:  config.TaskInfo_CompletedSubtaskCountField,
			taskmodel.SubtaskStatus_Timeout:   config.TaskInfo_TimeoutSubtaskCountField,
			taskmodel.SubtaskStatus_Cancelled: config.TaskInfo_CancelledSubtaskCountField,
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

func ReadSubtaskStatus(subtaskId uint64, statusRet *uint32) error {

	cmd := redistool.DefaultRedis().HGet(context.Background(), tasktool.GetSubtaskKey(subtaskId),
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

func IsSubtaskRunning(subtaskId uint64) bool {

	var status uint32 = 0
	err := ReadSubtaskStatus(subtaskId, &status)
	if err != nil {
		glog.Warning("failed to read subtask status: ", subtaskId, err)
		return false
	}

	return status == taskmodel.SubtaskStatus_Running
}
