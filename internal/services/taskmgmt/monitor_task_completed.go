package taskmgmt

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/dbdef"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/tasktool"
)

// 协程 <<go_monitor_completed_task>>
// 检查并处理已完成队列中的任务
func MonitorCompletedTask() {
	// 从已完成任务列表中取完成的任务id
	var taskList []taskmodel.TaskIdType
	getCompletedTask(&taskList)

	for _, taskId := range taskList {
		completeTask(taskId)
	}
}

// 从zset中取完成的任务ID
func getCompletedTask(taskList *[]taskmodel.TaskIdType) {

	// 取已完成列表中的最先完成的5个任务
	currentTime := strconv.FormatUint(uint64(time.Now().Unix()), 10)
	opt := redis.ZRangeBy{
		Min: "-inf", Max: currentTime,
		Offset: 0, Count: 5,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(context.Background(), config.CompletedTaskList, &opt)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get completed taskid: ", err.Error())
		return
	}

	// 解析出任务ID
	strArr, _ := cmd.Result()
	for _, str := range strArr {
		taskId, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert task id: ", str, err.Error())
			continue
		}

		*taskList = append(*taskList, taskmodel.TaskIdType(taskId))
	}

	glog.Info("succeeded to get completed tasks: ", *taskList)
}

// 完成任务
func completeTask(taskId taskmodel.TaskIdType) {

	// 从task list中删除任务记录，避免被monitor_task_timeout处理
	cmd := redistool.DefaultRedis().ZRem(context.Background(), config.TaskZset, taskId)
	_, err := cmd.Result()
	if err != nil {
		glog.Warning("failed to zrem task from task list: ", err.Error())
	}

	// 从已完成队列中删除任务记录
	cmd = redistool.DefaultRedis().ZRem(context.Background(), config.CompletedTaskList, taskId)
	val, err := cmd.Result()
	if err != nil {
		glog.Warning("failed to zrem task from completed task list: ", err.Error())
	}

	// 为0, 表示被其他例程处理了
	if val == 0 {
		glog.Info("owned by other: ", taskId)
		return
	}

	glog.Info("owned completed task: ", taskId)

	// 将取到的任务设置为已完成
	var taskRecord = dbdef.TaskRecord{}
	err = tasktool.CompleteTask(taskId, &taskRecord)
	if err != nil {
		glog.Warning("failed to complete task: ", taskId, err.Error())
		return
	}

	// 镜像安全类型，执行各类别完成回调
	AfterTaskCompleted(&taskRecord)

	// 执行清理操作
	cleanTaskKeys(taskId)

	glog.Info("succeeded to complete task: ", taskId)
}

// 清理任务的redis key
func cleanTaskKeys(taskId taskmodel.TaskIdType) error {

	pipeline := redistool.DefaultRedis().Pipeline()

	// 从redis_task_schedule_zset 中删除taskid
	pipeline.ZRem(context.Background(), config.GeneratingTaskZset, taskId)
	pipeline.ZRem(context.Background(), config.RunningTaskZset, taskId)
	pipeline.ZRem(context.Background(), config.ToGenerateTaskZset, taskId)

	// 执行pipeline
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec clean task keys pipeline: ", err)
		return err
	}

	return nil
}

// 任务完成时，执行任务类型的完成回调
func AfterTaskCompleted(taskRecord *dbdef.TaskRecord) error {
	return nil
}
