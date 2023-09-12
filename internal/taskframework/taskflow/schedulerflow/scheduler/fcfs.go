package scheduler

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/redistool"
)

// FCFS调度算法
type FCFS struct {
	QueueKeyName string // 队列的Key名
}

// FCFS调度算法
func (scheduler *FCFS) Schedule(
	retTaskId *taskmodel.TaskIdType,
	noTask *bool,
) error {
	return ScheduleQueue(scheduler.QueueKeyName, retTaskId, noTask)
}

// 从队列key中取任务ID
func ScheduleQueue(
	keyName string,
	retTaskId *taskmodel.TaskIdType,
	noTask *bool,
) error {

	// 按照FIFO的策略，取队首的元素
	err := PopTaskId(keyName, retTaskId)

	// 如果没有元素，返回
	if err == errordef.ErrNotFound {
		*noTask = true
		return nil
	}

	if err != nil {
		glog.Warning("failed to pop task id from ", keyName, ", ", err.Error())
		return err
	}

	return nil
}

// 取队首的任务ID元素
func PopTaskId(
	keyName string,
	retTaskId *taskmodel.TaskIdType,
) error {

	if len(keyName) == 0 {
		glog.Error("empty keyname")
		return errors.New("empty keyname")
	}

	// 取队列首部的元素
	cmd := redistool.DefaultRedis().LPop(context.Background(), keyName)
	err := cmd.Err()

	// 队列中没有元素
	if err == redis.Nil {
		return errordef.ErrNotFound
	}

	if err != nil {
		glog.Warning("failed to lpop task id: ", err.Error())
		return err
	}

	// 转换为任务ID
	taskId, err := cmd.Uint64()
	if err != nil {
		glog.Warning("failed to convert lpopped task id: ", err.Error())
		return err
	}

	*retTaskId = taskmodel.TaskIdType(taskId)
	return nil
}
