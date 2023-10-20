package scheduler

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
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

func MonitorRunningTaskToComplete() {
	// 检查是否有已完成的任务
	var taskList = []uint64{}
	err := getToBeCompletedTask(&taskList)
	if err != nil {
		glog.Warning("failed to get completed tasks: ", err)
		return
	}

	if len(taskList) <= 0 {
		return
	}

	// 对任务进行完成操作
	err = completeTask(&taskList)
	if err != nil {
		glog.Warning("failed to process completed tasks: ", taskList, err)
		return
	}
}

func getToBeCompletedTask(taskList *[]uint64) error {

	if taskList == nil {
		panic("invalid taskList pointer")
	}

	// 取redis_task_schedule_zset 中的 $taskid
	opt := redis.ZRangeBy{
		Min: "-inf", Max: "+inf",
		Offset: 0, Count: 10,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(
		context.Background(), config.RunningTaskZset, &opt,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get completed tasks from redis: ", err)
		return err
	}

	strList := cmd.Val()
	glog.Info("got running tasks: ", strList)

	// 转换查询到的任务ID
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert completed subtask id: ", str)
			continue
		}

		// 检查任务是否已经完成, 或者为过久的任务
		taskId := taskmodel.TaskIdType(id)
		if !tasktool.CheckIfTaskCompleted(taskId) && !checkIfTooLongTask(taskId) {
			continue
		}

		// 任务已经完成，添加到队列中
		*taskList = append(*taskList, id)
	}

	// 如果列表为空，表示没有完成的任务
	if len(*taskList) == 0 {
		glog.Info("get empty list, no completed task")
		return nil
	}

	glog.Info("got completed tasks: ", *taskList)
	return nil
}

// 检查任务是否创建了太久时间
func checkIfTooLongTask(taskId taskmodel.TaskIdType) bool {

	var createTime uint64 = 0
	err := tasktool.GetTaskCreateTime(taskId, &createTime)

	// 字段不存在表示任务已经创建很久，key已经失效
	if err == errordef.ErrNotFound {
		return true
	}

	// 其他错误，认为无法判断
	if err != nil {
		return false
	}

	return time.Now().Unix()-int64(createTime) > 24*3600
}

// 批量执行任务完成的操作
func completeTask(taskList *[]uint64) error {

	if taskList == nil {
		panic("invalid task list pointer")
	}

	// 批量从redis_task_schedule_zset中删除taskid
	var ownedTasks = []uint64{}
	err := redistool.OwnElementsInList(config.RunningTaskZset, taskList, &ownedTasks)
	if err != nil {
		return err
	}

	if len(ownedTasks) == 0 {
		return nil
	}

	// 处理有所有权的任务, 执行任务的完成操作
	pipeline := redistool.DefaultRedis().Pipeline()
	for _, taskId := range ownedTasks {
		err = PerformCompleteTask(taskId, &pipeline)
		if err != nil {
			glog.Warning("failed to complete task: ", taskId, err)
		}
	}

	// 执行pipeline
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	return nil
}

// 设置任务完成, 执行一些完成操作
func PerformCompleteTask(taskId uint64, ppipeline *redis.Pipeliner) error {

	pipeline := *ppipeline

	// 将任务id放到已完成列表中
	z := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: taskId,
	}
	pipeline.ZAdd(context.Background(), config.CompletedTaskList, &z)

	// 清理临时的redis key.

	return nil
}
