package scheduler

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/subtasktool"
)

func MonitorTimeoutSubtask() {
	// 取超时的子任务
	var subtaskList = []uint64{}
	err := getTimeoutSubtasks(&subtaskList)
	if err != nil {
		glog.Warning("failed to get timeout subtasks: ", err)
		return
	}

	if len(subtaskList) <= 0 {
		return
	}

	// 处理超时的子任务
	err = repairTimeoutSubtasks(&subtaskList)
	if err != nil {
		glog.Warning("failed to repair timeout subtask: ", err)
		return
	}
}

func getTimeoutSubtasks(subtaskList *[]uint64) error {

	if subtaskList == nil {
		panic("invalid subtaskList pointer")
	}

	// 从redis_subtask_scanning_zset 中取超时的子任务
	now := time.Now().Unix()
	nowStr := strconv.FormatUint(uint64(now), 10)
	opt := redis.ZRangeBy{
		Min: "-inf", Max: nowStr,
		Offset: 0, Count: 100,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(
		context.Background(), config.RunningSubtaskZset, &opt,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get timeout subtask from redis: ", err)
		return err
	}

	strList := cmd.Val()
	if len(strList) > 0 {
		glog.Info("got timeout subtask: ", strList)
	}

	// 转换查询到的子任务ID
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert timeout subtask id: ", str)
			continue
		}

		*subtaskList = append(*subtaskList, id)
	}

	// no timeout subtask
	if len(*subtaskList) == 0 {
		glog.Info("get empty list, no timeout subtask")
		return nil
	}

	glog.Info("got timeout subtasks: ", *subtaskList)
	return nil
}

// 修复子任务的超时状态
func repairTimeoutSubtasks(subtaskList *[]uint64) error {

	if subtaskList == nil {
		panic("invalid subtaskList pointer")
	}

	// remove this subtask from running subtasks list
	owndSubtaskList := []uint64{}
	err := redistool.OwnElementsInList(config.RunningSubtaskZset, subtaskList, &owndSubtaskList)
	if err != nil {
		return err
	}

	if len(owndSubtaskList) <= 0 {
		return nil
	}

	completeTime := time.Now().Unix()
	pipeline := redistool.DefaultRedis().Pipeline()
	for _, id := range owndSubtaskList {

		glog.Info("owned subtask, set subtask to complete: ", id)

		// set completion code to timeout
		err = subtasktool.SetSubtaskResult(id, taskmodel.SubtaskResult_Timeout, "", &pipeline)
		if err != nil {
			glog.Warning("failed to set subtask timeout: ", id, err)
		}

		// 插入到 redis_subtask_complete_list 完成队列中
		z := redis.Z{
			Member: id,
			Score:  float64(completeTime),
		}

		pipeline.ZAdd(context.Background(), config.CompletedSubtaskList, &z)
	}

	_, err = pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	glog.Info("succeeded to repair timeout subtasks: ", *subtaskList)
	return nil
}
