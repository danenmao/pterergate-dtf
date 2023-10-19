package tasktool

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

// 创建子任务信息key
func CreateSubtaskInfoKey(
	subtaskId uint64,
	subtaskData *taskmodel.SubtaskBody,
) error {

	var valueMap = map[string]interface{}{
		config.SubtaskInfo_UID:            0,
		config.SubtaskInfo_TaskIdField:    uint64(subtaskData.TaskId),
		config.SubtaskInfo_PriorityField:  0,
		config.SubtaskInfo_StartTimeField: time.Now().Unix(),
		config.SubtaskInfo_Param:          subtaskData.TypeParam,
		config.SubtaskInfo_StatusField:    taskmodel.SubtaskStatus_Running,
		config.SubtaskInfo_TaskTypeField:  subtaskData.TaskType,
	}

	// 设置子任务的运行信息
	cmd := redistool.DefaultRedis().HMSet(
		context.Background(),
		GetSubtaskKey(subtaskId),
		valueMap,
	)

	redistool.DefaultRedis().Expire(
		context.Background(),
		GetSubtaskKey(subtaskId),
		time.Hour*4,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to create subtask info key: ", subtaskId, ",", err)
		return err
	}

	return nil
}

// 将子任务添加到任务中
func AddSubtaskToTask(taskId taskmodel.TaskIdType, subtaskId uint64) error {

	pipeline := redistool.DefaultRedis().Pipeline()

	// 将子任务推入 redis_subtask_list.$taskid，zset，按生成时间排序
	pipeline.ZAdd(context.Background(), GetTaskSubtaskListKey(taskId), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: subtaskId,
	})

	// 修改 redis_task_info.$taskid中的subtaskcount
	pipeline.HIncrBy(context.Background(), GetTaskInfoKey(taskId),
		config.TaskInfo_TotalSubtaskCountField, 1)

	// 镜像任务的子任务不进入redis_subtask_to_schedule_list, 直接进入redis_subtask_scanning_zset

	// 执行pipeline
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", subtaskId, err)
		return err
	}

	glog.Info("succeeded to add subtask to task: ", subtaskId, taskId)
	return nil
}

// 将子任务添加到redis_subtask_scanning_zset, 执行中的子任务列表
func AddSubtaskToRunningList(
	subtasks *[]taskmodel.SubtaskBody,
) error {

	// 拼装添加命令
	const DefaultTimeout = 720
	zlist := []*redis.Z{}
	for _, subtask := range *subtasks {

		timeout := DefaultTimeout
		if subtask.Timeout != 0 {
			timeout = int(subtask.Timeout)
		}

		zlist = append(zlist, &redis.Z{
			Score:  float64(time.Now().Add(time.Duration(timeout) * time.Second).Unix()),
			Member: uint64(subtask.SubtaskId),
		})
	}

	// 将子任务推入 redis_subtask_scanning_zset, zset, 按照超时时间排序
	pipeline := redistool.DefaultRedis().Pipeline()
	pipeline.ZAdd(context.Background(), config.RunningSubtaskZset, zlist...)

	// 执行pipeline
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec AddSubtaskToRunningList pipeline: ", err)
		return err
	}

	return nil
}

func GetTaskIdOfSubtask(subtaskId uint64, taskId *taskmodel.TaskIdType) error {

	cmd := redistool.DefaultRedis().HGet(
		context.Background(), GetSubtaskKey(subtaskId), config.SubtaskInfo_TaskIdField,
	)

	err := cmd.Err()
	if err == redis.Nil {
		return errordef.ErrNotFound
	}

	if err != nil {
		glog.Warning("failed to read task id of subtask: ", subtaskId, ", ", err)
		return err
	}

	idStr := cmd.Val()
	var intId uint64 = 0
	intId, err = strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		glog.Warning("failed to convert task id str: ", idStr, err)
		return err
	}

	*taskId = taskmodel.TaskIdType(intId)
	return nil
}

//
