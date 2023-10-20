package executorconnector

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
)

// 重试推送到执行器服务的队列的名称
const RedisRetryToPushExecutorQueue = "retry.push.to.executor.queue"
const RetryToPushInterval = 2

type RetrySubtaskData struct {
	taskmodel.SubtaskBody
	ExpiredAt time.Time `json:"expired_at"` // 子任务重试的截止时间
}

func AddSubtasksToRetryQueue(
	subtasks *[]taskmodel.SubtaskBody,
) error {

	// 将子任务数据批量序列化
	now := time.Now()
	vals := []interface{}{}
	for _, subtask := range *subtasks {

		retryData := RetrySubtaskData{
			SubtaskBody: subtask,
			ExpiredAt:   subtask.CreatedAt.Add(time.Second * time.Duration(subtask.Timeout)),
		}

		// 剔除已经超时的子任务，不重试
		if retryData.ExpiredAt.Sub(now.Add(time.Second*2)) <= 0 {
			continue
		}

		data, err := json.Marshal(&retryData)
		if err != nil {
			glog.Warning("failed to marshal subtask data: ", subtask.SubtaskId, subtask.TaskId)
			continue
		}

		vals = append(vals, string(data))
	}

	// 批量保存到Redis重试队列中
	cmd := redistool.DefaultRedis().RPush(context.Background(), RedisRetryToPushExecutorQueue, vals...)
	redistool.DefaultRedis().Expire(context.Background(), RedisRetryToPushExecutorQueue, time.Hour*8)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add subtasks to retry queue")
		return err
	}

	glog.Info("succeeded to add subtasks to retry queue: ", len(*subtasks))
	return nil
}

// 重试例程
// 重试将子任务推送给执行器服务
func RetryPushToExecutorRoutine() {
	routine.ExecRoutineWithInterval(
		"RetryPushToExecutorRoutine",
		retryPushToExecutor,
		time.Duration(RetryToPushInterval)*time.Second,
	)
}

// 重试将子任务推送给执行器服务
func retryPushToExecutor() {

	glog.Info("retry to push subtasks to executor")

	// 取子任务列表
	subtasks := []taskmodel.SubtaskBody{}
	err := getRetryPushSubtasks(&subtasks)
	if err != nil {
		glog.Warning("failed to get subtasks to retry to push: ", err)
		return
	}

	if len(subtasks) <= 0 {
		glog.Info("no subtask to retry")
		return
	}

	// 推送给执行器
	// 将子任务分批发送给执行器服务
	failedSubtasks := []taskmodel.SubtaskBody{}
	err = PushToExecutor(&subtasks, &failedSubtasks)
	if err != nil {
		glog.Error("failed to push subtasks to executor: ", err.Error())
	}

	// 如果推送失败，将子任务放到失败重试队列中, 稍后重试
	if len(failedSubtasks) > 0 {
		AddSubtasksToRetryQueue(&failedSubtasks)
	}
}

// 取要重试的子任务列表
func getRetryPushSubtasks(
	subtasks *[]taskmodel.SubtaskBody,
) error {

	// 构造命令pipeline
	pipeline := redistool.DefaultRedis().Pipeline()
	for i := uint32(0); i < ExecutorMaxPushSubtaskCount; i++ {
		pipeline.LPop(context.Background(), RedisRetryToPushExecutorQueue)
	}

	// 执行
	cmdList, err := pipeline.Exec(context.Background())
	if err == redis.Nil {
		glog.Info("retry key not exist")
		return nil
	}

	if err != nil {
		glog.Warning("failed to exec pop retry subtask list pipeline: ", err)
		return err
	}

	// 从pipeline结果中读取子任务列表
	now := time.Now()
	for _, cmd := range cmdList {
		err = cmd.Err()

		// 无元素, 退出
		if err == redis.Nil {
			break
		}

		if err != nil {
			glog.Warning("retry pipeline cmd return error: ", err)
			break
		}

		strCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			glog.Error("invalid retry cmd type: ", cmd)
			continue
		}

		// 读取保存的子任务数据
		data := []byte(strCmd.Val())
		retryData := RetrySubtaskData{}
		err = json.Unmarshal(data, &retryData)
		if err != nil {
			glog.Warning("failed to unmarshal retry subtask: ", string(data), ",", err)
			continue
		}

		// 剔除已经超时的子任务，不重试
		if retryData.ExpiredAt.Sub(now.Add(time.Second*2)) <= 0 {
			glog.Info("remove timeout subtask in retry queue: ", retryData.TaskId, ",", retryData.SubtaskId)
			continue
		}

		*subtasks = append(*subtasks, retryData.SubtaskBody)
	} // for

	glog.Info("succeeded to get retry push subtasks: ", len(*subtasks))
	return nil
}
