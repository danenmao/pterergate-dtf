package executorconnector

import (
	"time"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
)

// 批量推送子任务的上限
const ExecutorMaxPushSubtaskCount uint32 = 10

var ExecutorService taskmodel.ExecutorInvoker

// 执行子任务列表
// 传入的子任务都属于同一个任务
func ExecSubtasks(
	taskId taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskData,
) error {

	// 将子任务分批发送给执行器服务, image_subtask_executor
	failedSubtasks := []taskmodel.SubtaskData{}
	err := PushToExecutor(subtasks, &failedSubtasks)
	if err != nil {
		glog.Error("failed to push subtasks to executor: ", err.Error())
	}

	// 如果推送失败，将子任务放到失败重试队列中, 稍后重试
	if len(failedSubtasks) > 0 {
		AddSubtasksToRetryQueue(&failedSubtasks)
	}

	return err
}

// 将子任务分批发送给执行器服务
func PushToExecutor(
	subtasks *[]taskmodel.SubtaskData,
	failedSubtasks *[]taskmodel.SubtaskData,
) error {

	// 将子任务发送给执行器服务
	totalCount := len(*subtasks)
	for i := 0; i < totalCount; i += int(ExecutorMaxPushSubtaskCount) {

		// 计算批量范围
		start := i
		end := i + int(ExecutorMaxPushSubtaskCount)
		if end > totalCount {
			end = totalCount
		}

		batchList := (*subtasks)[start:end]

		// 将子任务批量发送给执行器服务
		err := PushBatchSubtaskToExecutor(batchList, failedSubtasks)
		if err != nil {
			*failedSubtasks = append(*failedSubtasks, batchList...)
			glog.Info("added failed subtasks to retry queue: ", len(batchList))
		}

		time.Sleep(time.Millisecond)
	}

	return nil
}

// 将一批子任务发送给执行器服务
func PushBatchSubtaskToExecutor(
	subtasks []taskmodel.SubtaskData,
	retFailedSubtasks *[]taskmodel.SubtaskData,
) error {

	glog.Infof("ready to push batch subtask, subtask num: %d", len(subtasks))

	failedSubtasks := []taskmodel.SubtaskData{}
	err := sendRequestToExecutor(subtasks)
	if err != nil {
		return err
	}

	// 处理失败的子任务项
	glog.Infof("total: %d, failed: %d", len(subtasks), len(failedSubtasks))
	if len(failedSubtasks) > 0 {
		*retFailedSubtasks = append(*retFailedSubtasks, failedSubtasks...)
	}

	return nil
}

// 向执行器发送请求
func sendRequestToExecutor(
	req []taskmodel.SubtaskData,
) error {

	if len(req) <= 0 {
		return nil
	}

	return ExecutorService(req)
}
