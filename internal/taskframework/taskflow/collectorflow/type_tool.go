package collectorflow

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
	"github.com/danenmao/pterergate-dtf/internal/subtasktool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskloader"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

func GetTaskCollectorCallback(taskType uint32, collector *taskmodel.ITaskCollectorCallback) error {

	var plugin taskplugin.ITaskPlugin = nil
	err := taskloader.LookupTaskPlugin(taskType, &plugin)
	if err != nil {
		glog.Warning("failed to get task plugin: ", taskType)
		return err
	}

	var context taskmodel.TaskBody
	err = plugin.GetTaskBody(&context)
	if err != nil {
		glog.Warning("failed to get task context: ", err.Error())
		return err
	}

	*collector = context.CollectorCallback
	glog.Info("succeeded to get task collector: ", taskType)
	return nil
}

// 从子任务ID得到collector
func GetSubtaskCollectorCallback(
	subtaskId taskmodel.SubtaskIdType,
	collector *taskmodel.ITaskCollectorCallback,
) error {

	// 取子任务的任务类型
	taskType := uint32(0)
	err := subtasktool.GetSubtaskTaskType(uint64(subtaskId), &taskType)
	if err != nil {
		glog.Warning("failed to get task type of subtask: ", subtaskId, ",", err)
		return err
	}

	// 获取collector
	err = GetTaskCollectorCallback(taskType, collector)
	if err != nil {
		glog.Warning("failed to get task type collector: ", taskType, ",", err)
		return err
	}

	glog.Info("succeeded to get task type collector: ", taskType)
	return nil
}

// 从子任务ID得到collector
func GetCollectorCallbackByTaskId(
	taskId taskmodel.TaskIdType,
	collector *taskmodel.ITaskCollectorCallback,
) error {

	// 取任务的任务类型
	taskType := uint32(0)
	err := tasktool.GetTaskType(taskId, &taskType)
	if err != nil {
		glog.Warning("failed to get task type: ", taskId, ",", err)
		return err
	}

	err = GetTaskCollectorCallback(taskType, collector)
	if err != nil {
		glog.Warning("failed to get task type collector: ", taskType, ",", err)
		return err
	}

	glog.Info("succeeded to get task type collector: ", taskType)
	return nil
}
