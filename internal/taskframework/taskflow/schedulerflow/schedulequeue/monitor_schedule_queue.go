package schedulequeue

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/taskframework/taskflow/flowdef"
	"pterergate-dtf/internal/tasktool"
)

// 调度中任务队列
const CurrentTaskZSet = "current.task.list"
const DefaultCurrentTaskTimeout = time.Minute
const CheckCurrentTaskInterval = time.Duration(10) * time.Second

// 添加到调度中任务队列
func AddToCurrentTaskList(taskId taskmodel.TaskIdType) {

	timeout := time.Now().Add(DefaultCurrentTaskTimeout).Unix()
	cmd := redistool.DefaultRedis().ZAdd(context.Background(), CurrentTaskZSet, &redis.Z{
		Score:  float64(timeout),
		Member: taskId,
	})

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add task to current task list: ", taskId, ", err:", err)
	}

	glog.Info("added task to current task list: ", taskId)
}

// 添加到调度中任务队列
func AddListToCurrentTaskList(taskList []taskmodel.TaskIdType) {

	if len(taskList) <= 0 {
		return
	}

	timeout := time.Now().Add(DefaultCurrentTaskTimeout).Unix()
	zlist := []*redis.Z{}
	for _, taskId := range taskList {
		zlist = append(zlist, &redis.Z{
			Score:  float64(timeout),
			Member: taskId,
		})
	}

	cmd := redistool.DefaultRedis().ZAdd(context.Background(), CurrentTaskZSet, zlist...)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add task to current task list: ", taskList, ", err:", err)
	}
}

// 从当前任务列表中删除任务
func RemoveFromCurrentTaskList(taskId taskmodel.TaskIdType, pipeline redis.Pipeliner) {
	pipeline.ZRem(context.Background(), CurrentTaskZSet, taskId)
}

// 从当前任务列表中删除任务
func RemoveFromCurrentTaskListDirectly(taskId taskmodel.TaskIdType) {
	redistool.DefaultRedis().ZRem(context.Background(), CurrentTaskZSet, taskId)
}

// 从当前任务列表中删除任务
func RemoveListFromCurrentTaskList(taskIdList []interface{}, pipeline redis.Pipeliner) {
	pipeline.ZRem(context.Background(), CurrentTaskZSet, taskIdList...)
}

// 协程 <<go_monitor_current_task>>
func MonitorCurrentTaskRoutine() {
	routine.ExecRoutineByDuration(
		"MonitorCurrentTaskRoutine",
		monitorCurrentTask,
		CheckCurrentTaskInterval,
	)
}

// 监控调度中的当前任务
func monitorCurrentTask() {

	// 取已完成的子任务
	var taskList = []uint64{}
	err := getLostTask(&taskList)
	if err != nil {
		glog.Warning("failed to get lost tasks: ", err)
		return
	}

	if len(taskList) <= 0 {
		return
	}

	// 处理已完成的子任务
	err = repairLostTaskList(&taskList)
	if err != nil {
		glog.Warning("failed to process lost tasks: ", taskList, ", err: ", err)
		return
	}
}

// 获取在调度中丢失的任务
func getLostTask(retTaskList *[]uint64) error {

	if retTaskList == nil {
		panic("invalid taskList pointer")
	}

	// 从current_task_list中取超时的子任务
	taskList := []uint64{}
	err := redistool.GetTimeoutElemList(CurrentTaskZSet, 100, &taskList)
	if err != nil {
		glog.Warning("failed to get lost task from redis: ", err)
		return err
	}

	// 如果列表为空，表示没有丢失的任务
	if len(taskList) == 0 {
		return nil
	}

	ownedTaskList := []uint64{}
	err = redistool.TryToOwnElemList(CurrentTaskZSet, &taskList, &ownedTaskList)
	if err != nil {
		glog.Warning("failed to own lost task: ", taskList)
		return err
	}

	if len(ownedTaskList) <= 0 {
		return nil
	}

	*retTaskList = append(*retTaskList, ownedTaskList...)
	glog.Info("got lost tasks: ", *retTaskList)
	return nil
}

// 修复丢失中的任务
func repairLostTaskList(taskList *[]uint64) error {

	if taskList == nil {
		panic("invalid taskList pointer")
	}

	if len(*taskList) <= 0 {
		return nil
	}

	pipeline := redistool.DefaultRedis().Pipeline()
	for _, taskId := range *taskList {
		repairTask(taskId, pipeline)
	}

	// 执行pipeline
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec repiar lost task pipeline: ", taskList, ", err:", err)
		return err
	}

	return nil
}

// 修改丢失的任务
func repairTask(taskId uint64, pipeline redis.Pipeliner) error {

	// 读取任务的调度数据
	data := flowdef.TaskScheduleData{}
	err := tasktool.GetTaskScheduleData(taskmodel.TaskIdType(taskId), &data)
	if err != nil {
		glog.Warning("failed to get task schedule data while repair task: ", taskId, ",", err)
		return err
	}

	if len(data.CurrentQueueKeyName) <= 0 {
		glog.Warning("empty CurrentQueueKeyName: ", taskId)
		return errors.New("empty CurrentQueueKeyName")
	}

	// 将任务追加到原调度队列尾部
	cmd := pipeline.RPush(context.Background(), data.CurrentQueueKeyName, uint64(taskId))
	err = cmd.Err()
	if err != nil {
		glog.Warning("failed to append lost task to queue key: ", taskId, ", ", err)
		return err
	}

	glog.Info("succeeded to append lost task to queue key: ", taskId, ",", data.CurrentQueueKeyName)
	return nil
}
