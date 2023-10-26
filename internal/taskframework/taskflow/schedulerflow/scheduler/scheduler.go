package scheduler

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

type IQueueScheduler interface {
	Schedule(taskId *taskmodel.TaskIdType, noTask *bool) error
}

// 从队列中取出任务ID
func ScheduleTaskInQueue(
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
