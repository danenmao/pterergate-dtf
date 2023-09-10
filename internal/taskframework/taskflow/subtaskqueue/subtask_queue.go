package subtaskqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/redistool"
)

// 任务的子任务队列
const (
	// 任务下的待调度的子任务队列, to_schedule_subtask_list.$taskid, list
	ToScheduleSubtaskSetOfTaskPrefix = "dtf.task.to.schedule.subtask.list."
)

// 获取任务的子任务队列的键名
func GetSubtaskQueueOfTask(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", ToScheduleSubtaskSetOfTaskPrefix, taskId)
}

// 任务的子任务队列
// 一个任务一个子任务队列
type SubtaskQueue struct {
	TaskId taskmodel.TaskIdType
}

// 将子任务放入子任务队列中
func (queue *SubtaskQueue) PushSubtask(subtask *taskmodel.SubtaskData) error {

	// 序列化子任务
	data, err := json.Marshal(subtask)
	if err != nil {
		glog.Warning("failed to marshal subtask data: ", subtask.SubtaskId, ",", err)
		return err
	}

	// 将子任务放到任务的子任务队列中
	subtaskQueueKey := GetSubtaskQueueOfTask(taskmodel.TaskIdType(subtask.TaskId))
	cmd := redistool.DefaultRedis().RPush(context.Background(), subtaskQueueKey, string(data))
	redistool.DefaultRedis().Expire(context.Background(), subtaskQueueKey, time.Hour*8)

	err = cmd.Err()
	if err != nil {
		glog.Warning("failed to rpush subtask to task queue: ", subtask.SubtaskId, ",", err)
		return err
	}

	glog.Info("succeeded to rpush subtask to task queue: ", subtask.SubtaskId)
	return nil
}

// 从子任务队列中取子任务
func (queue *SubtaskQueue) PopSubtask(subtask *taskmodel.SubtaskData) error {

	// 从子任务队列中取子任务
	cmd := redistool.DefaultRedis().LPop(context.Background(),
		GetSubtaskQueueOfTask(taskmodel.TaskIdType(queue.TaskId)))
	err := cmd.Err()

	// 无子任务
	if err == redis.Nil {
		return errordef.ErrNotFound
	}

	// 其他错误
	if err != nil {
		glog.Warning("failed to pop subtask from queue: ", err)
		return err
	}

	// 反序列化出子任务数据
	err = json.Unmarshal([]byte(cmd.Val()), subtask)
	if err != nil {
		glog.Warning("failed to unmarshal subtask: ", err, ",", cmd.Val())
		return err
	}

	glog.Info("succeeded to pop subtask: ", subtask.TaskId)
	return nil
}

// 判断子任务队列中子任务的数量
func (queue *SubtaskQueue) GetSubtaskCount(taskId taskmodel.TaskIdType) (uint, error) {
	return GetSubtaskCount(taskId)
}

// 判断子任务队列中子任务的数量
func GetSubtaskCount(taskId taskmodel.TaskIdType) (uint, error) {

	cmd := redistool.DefaultRedis().LLen(context.Background(), GetSubtaskQueueOfTask(taskId))
	err := cmd.Err()

	// 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0
	// 如果 key 不是列表类型，返回一个错误。
	if err != nil {
		glog.Warning("failed to get subtask count from queue: ", taskId, ",", err)
		return 0, err
	}

	subtaskCount := uint(cmd.Val())
	glog.Info("get subtask count from queue: ", taskId, ", ", subtaskCount)
	return subtaskCount, nil
}
