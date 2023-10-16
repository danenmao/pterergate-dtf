package collectorflow

import (
	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
)

func OnSubtaskResult(
	subtaskResult *taskmodel.SubtaskResult,
	retFinished *bool,
) error {

	var collector taskmodel.ITaskCollector
	err := GetSubtaskCollector(subtaskResult.SubtaskId, &collector)
	if err != nil {
		glog.Warning("failed to get subtask collector: ", subtaskResult.SubtaskId, ",", err)
		return err
	}

	support := TaskCollectorSupport{}
	finished, err := collector.AfterExecution(subtaskResult, &support)
	if err != nil {
		glog.Warning("task type collector.OnScanResult return err: ", subtaskResult.SubtaskId, ",", err)
	}

	*retFinished = finished
	return nil
}

func OnSubtaskCompleted(
	subtaskResult *taskmodel.SubtaskResult,
) error {

	var collector taskmodel.ITaskCollector
	err := GetSubtaskCollector(subtaskResult.SubtaskId, &collector)
	if err != nil {
		glog.Warning("failed to get subtask collector: ", subtaskResult.SubtaskId, ",", err)
		return err
	}

	_, err = collector.AfterExecution(subtaskResult, nil)
	if err != nil {
		glog.Warning("task type collector.OnSubtaskCompleted return err: ", subtaskResult.SubtaskId, ",", err)
	} else {
		glog.Info("task type collector.OnSubtaskCompleted return code: ", subtaskResult.SubtaskId)
	}

	return nil
}

func AfterTaskCompleted(taskId taskmodel.TaskIdType) error {

	var collector taskmodel.ITaskCollector
	err := GetCollectorFromTaskId(taskId, &collector)
	if err != nil {
		glog.Warning("failed to get subtask collector ", taskId, ",", err)
		return err
	}

	code, err := collector.AfterTaskCompleted(taskId)
	if err != nil {
		glog.Warning("task type collector.AfterTaskCompleted return err: ", taskId, ",", err)
	} else {
		glog.Info("task type collector.AfterTaskCompleted return code: ", taskId, ",", code)
	}

	return nil
}
