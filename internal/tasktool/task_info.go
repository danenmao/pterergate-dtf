package tasktool

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/dbdef"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/subtaskqueue"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/tasklogicdef"
)

// 为任务创建info key
func CreateTaskInfoKey(
	taskId taskmodel.TaskIdType,
	taskParam *taskmodel.TaskParam,
) error {

	var data = map[string]interface{}{
		config.TaskInfo_UID:                        taskParam.Creator.UID,
		config.TaskInfo_StageField:                 config.Stage_CreatingTask,
		config.TaskInfo_StepField:                  1,
		config.TaskInfo_CreateTimeField:            time.Now().Unix(),
		config.SubtaskInfo_EndTimeField:            0,
		config.TaskInfo_TotalSubtaskCountField:     0,
		config.TaskInfo_CompletedSubtaskCountField: 0,
		config.TaskInfo_TimeoutSubtaskCountField:   0,
		config.TaskInfo_CancelledSubtaskCountField: 0,
		config.TaskInfo_GenerationCompletedField:   0,
		config.TaskInfo_ResourceCostField:          0,
		config.TaskInfo_TaskTypeField:              taskParam.TaskType,
		config.TaskInfo_Progess:                    0,
		config.TaskInfo_StatusField:                taskmodel.TaskStatus_Running,
	}

	taskKey := GetTaskInfoKey(taskId)
	cmd := redistool.DefaultRedis().HMSet(context.Background(), taskKey, data)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to create task info key: ", taskId, err.Error())
		return cmd.Err()
	}

	redistool.DefaultRedis().Expire(context.Background(), taskKey, time.Hour*72)

	glog.Info("succeeded to create task info key: ", taskKey)
	return nil
}

// 在Redis task-info key中保存任务的创建参数
func SetTaskRawTypeParam(taskId taskmodel.TaskIdType, paramStr string) error {

	taskInfoKey := GetTaskInfoKey(taskId)
	cmd := redistool.DefaultRedis().HSet(context.Background(), taskInfoKey, config.TaskInfo_TypeParam, paramStr)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to set task scan param: ", taskId, ", ", err)
		return err
	}

	glog.Info("succeeded to set task scan param: ", taskId)
	return nil
}

// 读取TypenParam
func GetTaskRawTypeParam(taskId taskmodel.TaskIdType, typeParam *string) error {

	if taskId == 0 {
		glog.Error("invalid task id: ", taskId)
		return errors.New("invalid task id")
	}

	taskInfoKey := GetTaskInfoKey(taskId)
	cmd := redistool.DefaultRedis().HGet(context.Background(), taskInfoKey, config.TaskInfo_TypeParam)
	if cmd.Err() != nil {
		glog.Warning("failed to get type param from task info key: ", taskId, cmd.Err())
		return cmd.Err()
	}

	data, err := cmd.Bytes()
	if err != nil {
		glog.Warning("failed to get data from redis cmd.Bytes: ", err)
		return err
	}

	*typeParam = string(data)
	return nil
}

// 获取任务的创建参数
func GetTaskCreateParam(
	taskId taskmodel.TaskIdType,
	retParam *tasklogicdef.TaskCreateParam,
) error {

	keyName := GetTaskCreateParamKey(taskId)
	cmd := redistool.DefaultRedis().Get(context.Background(), keyName)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get task create param key data: ", taskId, ", ", err.Error())
		return err
	}

	data := cmd.Val()
	err = json.Unmarshal([]byte(data), retParam)
	if err != nil {
		glog.Warning("failed to unmarshal create param: ", taskId, ", ", data)
		return err
	}

	return nil
}

// 保存任务的创建参数
func SaveTaskCreateParam(
	taskId taskmodel.TaskIdType,
	createParam *tasklogicdef.TaskCreateParam,
) error {

	data, err := json.Marshal(createParam)
	if err != nil {
		glog.Warning("failed to marshal create param: ", taskId, ", ", err.Error())
		return err
	}

	cmd := redistool.DefaultRedis().Set(context.Background(),
		GetTaskCreateParamKey(taskId),
		data,
		time.Hour*48,
	)
	err = cmd.Err()
	if err != nil {
		glog.Warning("failed to set task create param key: ", taskId, ", ", err.Error())
		return err
	}

	return nil
}

func GetInitTaskRecord(
	taskId taskmodel.TaskIdType,
	taskRecord *dbdef.DBTaskRecord,
) error {

	if taskId == 0 {
		glog.Warning("invalid task id")
		return errors.New("invalid task id")
	}

	taskInfoKey := GetTaskInfoKey(taskId)
	cmd := redistool.DefaultRedis().HGet(context.Background(), taskInfoKey, config.TaskInfo_InitTaskRecord)
	if cmd.Err() != nil {
		glog.Warning("failed to get init task record for task: ", taskId, cmd.Err())
		return cmd.Err()
	}

	err := json.Unmarshal([]byte(cmd.Val()), taskRecord)
	if err != nil {
		glog.Warning("failed to unmarshal init task record of task: ", taskId, err)
		return err
	}

	glog.Info("succeeded to get init task record of task: ", taskId)
	return nil
}

// 获取任务的类型
func GetTaskType(
	taskId taskmodel.TaskIdType,
	retTaskType *uint32,
) error {

	cmd := redistool.DefaultRedis().HGet(context.Background(), GetTaskInfoKey(taskId),
		config.TaskInfo_TaskTypeField)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get type of task: ", taskId, ",", err)
		return err
	}

	val := cmd.Val()

	taskType, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		glog.Warning("failed to parse task type field: ", taskId, ",", val, ",", err)
		return err
	}

	*retTaskType = uint32(taskType)
	return nil
}

// update next check time
func UpdateTaskGenerationNextCheckTime(taskId taskmodel.TaskIdType) error {

	cmd := redistool.DefaultRedis().HSet(
		context.Background(),
		GetTaskGenerationProgressKey(taskId),
		config.TaskGenerationKey_NextCheckTimeField,
		time.Now().Add(time.Minute).Unix(),
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to refresh task generation next check time value: ", taskId, err)
		return err
	}

	return nil
}

// 检查任务的生成是否完成
func CheckIfTaskGenerationCompleted(taskId taskmodel.TaskIdType) bool {

	// 检查redis_task_info.$taskid, 判断任务是否生成结束
	taskInfoKey := GetTaskInfoKey(taskId)
	strCmd := redistool.DefaultRedis().HGet(context.Background(),
		taskInfoKey, config.TaskInfo_GenerationCompletedField)
	err := strCmd.Err()
	if err != nil {
		glog.Warning("failed to get generation completed key of task: ", taskId, err)
		return false
	}

	completedStr := strCmd.Val()
	generationCompleted, err := strconv.Atoi(completedStr)
	if err != nil {
		glog.Warning("failed to convert generation completed key of task: ", taskId, completedStr, err)
		return false
	}

	// 生成尚未完成, 返回任务还未完成
	if generationCompleted == 0 {
		return false
	}

	glog.Info("task generation completed: ", taskId)
	return true
}

// 检查是否任务下的所有子任务都已经完成
func CheckIfAllSubtaskCompleted(taskId taskmodel.TaskIdType) bool {

	// redis_subtask_list.$taskid 为空
	subtaskListKey := GetTaskSubtaskListKey(taskId)
	cmd := redistool.DefaultRedis().ZCard(context.Background(), subtaskListKey)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get subtask count of task: ", taskId, err)
		return false
	}

	// 当key不存在时，值为0;
	// 当key为空时，值为0;
	count := cmd.Val()
	glog.Info("get subtask count of task: ", taskId, count)

	return count == 0
}

// 检查任务是否已经完成
func CheckIfTaskCompleted(taskId taskmodel.TaskIdType) bool {
	if taskId == 0 {
		glog.Warning("invalid task id: ", taskId)
		return false
	}

	// 检查redis_task_info.$taskid, 判断任务是否生成结束
	generationCompleted := CheckIfTaskGenerationCompleted(taskId)
	if !generationCompleted {
		return false
	}

	// 是否所有子任务都完成
	subtaskCompleted := CheckIfAllSubtaskCompleted(taskId)
	if !subtaskCompleted {
		return false
	}

	// 检查任务的本地子任务列表是否为空
	localSubtaskListEmpty := CheckIfLocalSubtaskListEmpty(taskId)
	return localSubtaskListEmpty
}

// 检查任务的本地子任务列表是否为空
func CheckIfLocalSubtaskListEmpty(taskId taskmodel.TaskIdType) bool {

	subtaskCount, err := subtaskqueue.GetSubtaskCount(taskmodel.TaskIdType(taskId))
	if err != nil {
		glog.Warning("failed to get local subtask count: ", taskId, ",", err)
		return false
	}

	return subtaskCount == 0
}

// 获取任务的优先级
func GetTaskPriority(
	taskId taskmodel.TaskIdType,
	retPriority *uint32,
) error {

	return nil
}

// 获取任务的调度数据
func GetTaskScheduleData(
	taskId taskmodel.TaskIdType,
	retScheduleData *tasklogicdef.TaskScheduleData,
) error {

	keyName := GetTaskScheduleDataKey(taskId)
	cmd := redistool.DefaultRedis().Get(context.Background(), keyName)
	err := cmd.Err()
	if err == redis.Nil {
		return &errordef.NotFoundError{}
	}

	if err != nil {
		glog.Warning("failed to get task schedule key data: ", taskId, ", ", err.Error())
		return err
	}

	data := cmd.Val()
	err = json.Unmarshal([]byte(data), retScheduleData)
	if err != nil {
		glog.Warning("failed to unmarshal schedule data: ", taskId, ", ", data)
		return err
	}

	return nil
}

// 保存任务的调度数据
func SaveTaskScheduleData(
	taskId taskmodel.TaskIdType,
	scheduleData *tasklogicdef.TaskScheduleData,
) error {

	data, err := json.Marshal(scheduleData)
	if err != nil {
		glog.Warning("failed to marshal schedule data: ", taskId, ", ", err.Error())
		return err
	}

	cmd := redistool.DefaultRedis().Set(
		context.Background(),
		GetTaskScheduleDataKey(taskId),
		data,
		time.Hour*48,
	)
	err = cmd.Err()
	if err != nil {
		glog.Warning("failed to set task schedule data key: ", taskId, ", ", err.Error())
		return err
	}

	return nil
}

func ReadSubtaskStartTime(
	subtaskId uint64,
	startTime *uint64, taskId *taskmodel.TaskIdType, uid *uint32,
	taskType *uint32,
) error {

	cmd := redistool.DefaultRedis().HGetAll(context.Background(), GetSubtaskKey(subtaskId))
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get info of subtask: ", subtaskId, err)
		return err
	}

	valMap := cmd.Val()

	// start_time
	timeStr, ok := valMap[config.SubtaskInfo_StartTimeField]
	if !ok {
		glog.Warning("no start_time field: ", subtaskId)
		return errors.New("no start_time field")
	}

	*startTime, err = strconv.ParseUint(timeStr, 10, 64)
	if err != nil {
		glog.Warning("failed to parse start_time field: ", subtaskId, timeStr, err)
		return err
	}

	// task id
	taskIdStr, ok := valMap[config.SubtaskInfo_TaskIdField]
	if !ok {
		glog.Warning("no task_id field: ", subtaskId)
		return errors.New("no task_id field")
	}

	var intId uint64 = 0
	intId, err = strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		glog.Warning("failed to convert task id: ", taskIdStr, err)
		return errors.New("failed to convert task id")
	}
	*taskId = taskmodel.TaskIdType(intId)

	// uid
	uidStr, ok := valMap[config.SubtaskInfo_UID]
	if !ok {
		glog.Warning("no owner field: ", subtaskId)
		return errors.New("no owner field")
	}

	uidValue, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		glog.Warning("failed to convert uid: ", uidStr, err)
		return errors.New("failed to convert uid")
	}

	*uid = uint32(uidValue)

	// task type
	taskTypeStr, ok := valMap[config.SubtaskInfo_TaskTypeField]
	if !ok {
		glog.Warning("no task type field found: ", subtaskId)
		taskTypeStr = "1"
	}

	taskTypeValue, err := strconv.ParseUint(taskTypeStr, 10, 32)
	if err != nil {
		glog.Warning("failed to convert task type: ", taskTypeStr, err)
		return errors.New("failed to convert task type")
	}

	*taskType = uint32(taskTypeValue)
	return nil
}

func GetTaskCreateTime(taskId taskmodel.TaskIdType, retCreateTime *uint64) error {

	cmd := redistool.DefaultRedis().HGet(context.Background(), GetTaskInfoKey(taskId),
		config.TaskInfo_CreateTimeField)
	err := cmd.Err()

	// 如果给定的字段或 key 不存在时，返回 nil
	if err == redis.Nil {
		return errordef.ErrNotFound
	}

	if err != nil {
		glog.Warning("failed to get create time of task: ", taskId, ", ", err)
		return err
	}

	timeStr := cmd.Val()
	*retCreateTime, err = strconv.ParseUint(timeStr, 10, 64)
	if err != nil {
		glog.Warning("failed to parse create_time field: ", taskId, ",", timeStr, ",", err)
		return err
	}

	return nil
}

func IsTaskRunning(taskId taskmodel.TaskIdType) bool {

	var status taskmodel.TaskStatusType = 0
	err := ReadTaskStatus(taskId, &status)
	if err != nil {
		glog.Warning("failed to read task status: ", taskId, err)
		return false
	}

	return status == taskmodel.TaskStatus_Running
}

func ReadTaskStatus(taskId taskmodel.TaskIdType, statusRet *taskmodel.TaskStatusType) error {

	cmd := redistool.DefaultRedis().HGet(context.Background(), GetTaskInfoKey(taskId),
		config.TaskInfo_StatusField)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get status of task: ", taskId, err)
		return err
	}

	val := cmd.Val()

	status, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		glog.Warning("failed to parse status field: ", taskId, val, err)
		return err
	}

	*statusRet = taskmodel.TaskStatusType(status)
	return nil
}
