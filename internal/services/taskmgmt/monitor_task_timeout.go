package taskmgmt

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

// 协程 <<go_monitor_task_timeout >>
// 用于检查执行超时的任务
func MonitorTaskTimeout() {

	// 从 redis_task_zset 中取超时的元素
	var taskList []taskmodel.TaskIdType
	getTimeoutTask(&taskList)

	// 设置这些任务为超时状态
	for _, taskId := range taskList {

		cmd := redistool.DefaultRedis().ZRem(context.Background(), config.TaskZset, taskId)
		val, err := cmd.Result()
		if err != nil {
			glog.Warning("failed to zrem task from task list: ", err.Error())
			continue
		}

		// 为0, 表示被其他例程处理了
		if val == 0 {
			glog.Info("owned by other: ", taskId)
			continue
		}

		err = setTaskTimeout(taskId)
		if err != nil {
			glog.Warning("failed to set task timeout: ", taskId)
			continue
		}

		glog.Info("succeeded to set task timeout: ", taskId)
	}
}

// 从zset中取超时的任务ID
func getTimeoutTask(taskList *[]taskmodel.TaskIdType) {

	// 检查过期时间戳在当前时间之前的元素数目
	currentTime := strconv.FormatUint(uint64(time.Now().Unix()), 10)
	countCmd := redistool.DefaultRedis().ZCount(context.Background(), config.TaskZset, "-inf", currentTime)
	_, err := countCmd.Result()
	if err != nil {
		glog.Warning("failed to get count of timeout task: ", err.Error())
		return
	}

	glog.Info("found timeout task, count: ", countCmd.Val())
	if countCmd.Val() == 0 {
		return
	}

	// 取前五个过期的元素
	opt := redis.ZRangeBy{
		Min: "-inf", Max: currentTime,
		Offset: 0, Count: 5,
	}

	rangeCmd := redistool.DefaultRedis().ZRangeByScore(context.Background(), config.TaskZset, &opt)
	err = rangeCmd.Err()
	if err != nil {
		glog.Warning("failed to get timeout taskid: ", err.Error())
		return
	}

	strArr, _ := rangeCmd.Result()
	for _, str := range strArr {
		taskId, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert task id: ", str, err.Error())
			continue
		}

		*taskList = append(*taskList, taskmodel.TaskIdType(taskId))
	}

	if len(*taskList) > 0 {
		glog.Info("get timeout task list: ", *taskList)
	}
}

// 设置任务已超时, 将任务推送到完成任务列表
func setTaskTimeout(taskId taskmodel.TaskIdType) error {
	return tasktool.PushTaskToCompletedList(taskId)
}

// 从运行中任务列表中删除任务
func RemoveFromRunningList(taskId taskmodel.TaskIdType) error {

	// 从redis_task_schedule_zset 中删除taskid
	err := redistool.DefaultRedis().ZRem(context.Background(), config.RunningTaskZset, taskId)
	if err != nil {
		glog.Warningf("failed to remove from running task list: %d, %s", taskId, err.Err().Error())
	}

	return nil
}
