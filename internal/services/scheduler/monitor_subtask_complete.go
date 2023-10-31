package scheduler

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

// 监视已完成的子任务
func MonitorSubtaskComplete() {
	// 取已完成的子任务
	var subtaskList = []uint64{}
	err := getCompletedSubtask(&subtaskList)
	if err != nil {
		glog.Warning("failed to get completed subtasks: ", err)
		return
	}

	if len(subtaskList) <= 0 {
		return
	}

	// 处理已完成的子任务
	err = processCompletedSubtask(&subtaskList)
	if err != nil {
		glog.Warning("failed to process completed subtasks: ", err)
		return
	}
}

// 获取已完成的子任务列表
func getCompletedSubtask(subtaskList *[]uint64) error {

	if subtaskList == nil {
		panic("invalid subtaskList pointer")
	}

	// 从subtask_complete_list中取完成的子任务
	opt := redis.ZRangeBy{
		Min: "-inf", Max: "+inf",
		Offset: 0, Count: 100,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(
		context.Background(),
		config.CompletedSubtaskList, &opt,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get completed subtask from redis: ", err)
		return err
	}

	strList := cmd.Val()
	if len(strList) > 0 {
		glog.Info("got completed subtask: ", strList)
	}
	// 转换查询到的子任务ID
	var wrongFormatList = []interface{}{}
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert completed subtask id: ", str)

			// 如果转换失败，说明数据格式错误，移除元素
			wrongFormatList = append(wrongFormatList, str)
			continue
		}

		*subtaskList = append(*subtaskList, id)
	}

	// 删除转换失败的子任务数据
	redistool.DefaultRedis().ZRem(context.Background(), config.CompletedSubtaskList, wrongFormatList...)

	// 如果列表为空，表示没有超时的任务
	if len(*subtaskList) == 0 {
		return nil
	}

	glog.Info("got completed subtasks: ", *subtaskList)
	return nil
}

// 处理已完成的子任务列表
func processCompletedSubtask(subtaskList *[]uint64) error {

	if subtaskList == nil {
		panic("invalid subtaskList pointer")
	}

	ownedSubtaskList := []uint64{}
	err := ownCompletedSubtasks(subtaskList, &ownedSubtaskList)
	if err != nil {
		return err
	}

	if len(ownedSubtaskList) <= 0 {
		return nil
	}

	pipeline := redistool.DefaultRedis().Pipeline()
	for _, subtaskId := range ownedSubtaskList {

		// 获取子任务所属的任务id
		var taskId taskmodel.TaskIdType = 0
		err := tasktool.GetTaskIdOfSubtask(subtaskId, &taskId)
		if err != nil {
			glog.Warning("failed to get task id of subtask: ", subtaskId, ", ", err)
			continue
		}

		// 从redis_subtask_list.$taskid 中删除子任务.
		pipeline.ZRem(context.Background(), tasktool.GetTaskSubtaskListKey(taskId), subtaskId)

		// 执行子任务后处理
		OnSubtaskCompleted(taskId, subtaskId)
	}

	// 执行pipeline
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	return nil
}

// 试图获取完成子任务的所有权
func ownCompletedSubtasks(subtaskList *[]uint64, ownedSubtaskList *[]uint64) error {
	return redistool.TryToOwnElements(config.CompletedSubtaskList, subtaskList, ownedSubtaskList)
}

// 执行子任务后处理
func OnSubtaskCompleted(taskId taskmodel.TaskIdType, subtaskId uint64) {

}
