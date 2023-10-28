package generator

import (
	"context"
	"math/rand"
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

// monitor_task_generation
func MonitorTaskGenerationRoutine() {
	// 检查当前实例生成的任务数是否超过上限
	if IsFull() {
		glog.Info("exceed generating limit")
		return
	}

	// 检查生成过程中异常的任务id
	var taskId taskmodel.TaskIdType = 0
	err := getExceptionalGenerationTask(&taskId)
	if err == errordef.ErrNotFound {
		return
	}

	if err != nil {
		glog.Warning("failed to get generating task: ", err)
		return
	}

	// 修复任务
	err = repairTaskGeneration(taskId)
	if err != nil {
		glog.Warning("failed to repair task generation: ", taskId, err)
	}
}

// 获取生成异常的任务
func getExceptionalGenerationTask(taskId *taskmodel.TaskIdType) error {

	// 取正在生成的任务id列表
	var taskList = []taskmodel.TaskIdType{}
	err := getGeneratingTaskList(&taskList)
	if err != nil {
		return err
	}

	// 检查列表中的id, 取异常的任务
	for _, id := range taskList {
		exceptional, err := isTaskGenerationExceptional(id)
		if err != nil {
			continue
		}

		if !exceptional {
			continue
		}

		// 尝试增加当前实例生成的任务数
		if !IncrIfNotFull() {
			glog.Info("exceed generating limit")
			continue
		}

		glog.Info("found an exceptional task: ", id)
		err = tasktool.TryToOwnTask(id)
		if err != nil {
			Decr()
			glog.Info("failed to own task: ", id, err)
			continue
		}

		*taskId = id
		return nil
	}

	return errordef.ErrNotFound
}

// 获取正在生成的任务列表
func getGeneratingTaskList(taskList *[]taskmodel.TaskIdType) error {

	// 取redis_task_generation_zset的元素数目
	intCmd := redistool.DefaultRedis().ZCard(context.Background(), config.GeneratingTaskZset)
	err := intCmd.Err()
	if err != nil {
		glog.Warning("failed to get zcard:", err)
		return err
	}

	var limit int64 = 10
	var offset int64 = 0
	zcard := intCmd.Val()
	if zcard > limit {
		offset = int64(rand.Intn(int(zcard - limit)))
	}

	// 按照时间从redis_task_generation_zset的随机位置取$taskid
	opt := redis.ZRangeBy{
		Min: "-inf", Max: "+inf",
		Offset: offset, Count: limit,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(
		context.Background(), config.GeneratingTaskZset, &opt,
	)

	err = cmd.Err()
	if err != nil {
		glog.Warning("failed to get generating task from redis: ", err)
		return err
	}

	// 如果列表为空，表示没有生成中的任务
	taskStrList := cmd.Val()
	if len(taskStrList) == 0 {
		glog.Info("get empty list, no generating task")
		return nil
	}

	glog.Info("got task: ", taskStrList)

	for _, taskIdStr := range taskStrList {
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			glog.Warning("failed to convert task id: ", taskIdStr)
			continue
		}

		*taskList = append(*taskList, taskmodel.TaskIdType(taskId))
	}

	return nil
}

// 检查任务生成过程是否异常
func isTaskGenerationExceptional(taskId taskmodel.TaskIdType) (bool, error) {
	// 检查redis_task_generation.$taskid.progress
	cmd := redistool.DefaultRedis().HGetAll(context.Background(), tasktool.GetTaskGenerationProgressKey(taskId))
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get generation progress key: ", taskId, err)
		return false, err
	}

	// 如果 redis_task_generation.$taskid.progress 不存在, 或者 next_check_time 过期,即状态异常
	valMap := cmd.Val()
	nextCheckTimeStr, ok := valMap[config.TaskGenerationKey_NextCheckTimeField]
	if !ok {
		glog.Info("no next_check_time key: ", taskId)
		return true, nil
	}

	nextCheckTime, err := strconv.ParseUint(nextCheckTimeStr, 10, 64)
	if err != nil {
		glog.Warning("failed to convert next_check_time: ", taskId, nextCheckTimeStr, err)
		return false, err
	}

	return nextCheckTime < uint64(time.Now().Unix()), nil
}

// 恢复任务的生成过程
func repairTaskGeneration(taskId taskmodel.TaskIdType) error {

	// 重新触发协程 go_task_generation, 恢复任务生成操作
	go TaskGenerationRoutine(taskId, true)
	time.Sleep(time.Second)

	return nil
}
