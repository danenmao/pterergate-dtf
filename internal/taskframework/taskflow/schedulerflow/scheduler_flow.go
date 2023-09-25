package schedulerflow

import (
	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/dtf/taskplugin"
	"pterergate-dtf/internal/taskframework/taskflow/schedulerflow/executorconnector"
	"pterergate-dtf/internal/taskframework/taskflow/schedulerflow/resourcegroup"
	"pterergate-dtf/internal/taskframework/taskflow/subtaskqueue"
	"pterergate-dtf/internal/taskframework/taskloader"
	"pterergate-dtf/internal/tasktool"
)

func AddTaskToScheduler(
	taskId taskmodel.TaskIdType,
	groupName string,
	taskType uint32,
	priority uint32,
) error {
	return resourcegroup.GetResourceGroupMgr().AddTask(groupName, taskId, taskType, priority)
}

// if no task, retTaskId is 0, subtasks is empty
func GetSubtaskToSchedule(
	retTaskId *taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskData,
) error {

	err := resourcegroup.GetResourceGroupMgr().Select(retTaskId, subtasks)
	if err != nil {
		glog.Warning("failed to select subtasks: ", err)
		return err
	}

	return nil
}

// execute subtasks belonging to the same task
func ExecSubtasks(
	taskId taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskData,
	toPushbackSubtask *[]taskmodel.SubtaskData,
) error {

	if len(*subtasks) <= 0 {
		return nil
	}

	// get the task type
	var taskType uint32 = 0
	err := tasktool.GetTaskType(taskId, &taskType)
	if err != nil {
		glog.Warning("failed to get task type: ", taskId, ", ", err.Error())
		return err
	}

	// get the task scheduler
	var scheduler taskmodel.ITaskScheduler = nil
	err = GetTaskScheduler(taskType, &scheduler)
	if err != nil {
		glog.Warning("failed to get task scheduler: ", taskType, ", ", err.Error())
		return err
	}

	glog.Info("get task type plugin scheduler: ", taskType, ", ", scheduler)

	// to dispatch subtasks
	doneSubtaskList := []taskmodel.SubtaskData{}
	for _, subtask := range *subtasks {
		err = DispatchSubtask(taskType, scheduler, &subtask)
		if err != nil {
			*toPushbackSubtask = append(*toPushbackSubtask, subtask)
			continue
		}

		doneSubtaskList = append(doneSubtaskList, subtask)
	}

	// to monitor these subtasks' running statuses
	tasktool.AddSubtaskToRunningList(&doneSubtaskList)

	err = executorconnector.ExecSubtasks(taskId, subtasks)
	if err != nil {
		glog.Error("failed to execute subtasks: ", taskId, ", ", err.Error())
	}

	return nil
}

func GetTaskScheduler(taskType uint32, scheduler *taskmodel.ITaskScheduler) error {

	var plugin taskplugin.ITaskPlugin = nil
	err := taskloader.LookupTaskPlugin(taskType, &plugin)
	if err != nil {
		glog.Warning("failed to get task plugin: ", taskType)
		return err
	}

	var taskBody taskmodel.TaskBody
	err = plugin.GetTaskBody(&taskBody)
	if err != nil {
		glog.Warning("failed to get task context: ", err.Error())
		return err
	}

	*scheduler = taskBody.Scheduler
	glog.Info("succeeded to get task scheduler: ", taskType)
	return nil
}

// 对子任务执行调度操作
func DispatchSubtask(
	taskType uint32,
	scheduler taskmodel.ITaskScheduler,
	subtask *taskmodel.SubtaskData,
) error {

	glog.Info("dipatch subtask: ", subtask)

	// invoke the dispatch method
	toDipatch, err := scheduler.BeforeDispatch(subtask.SubtaskId, subtask)
	if err != nil {
		glog.Info("dispatch subtask error,  pushed back: ", subtask, ", ", err.Error())
		return err
	}

	// don't dispatch it now, push it back
	if !toDipatch {
		subtaskqueue.PushSubtaskBack(subtask.TaskId, &[]taskmodel.SubtaskData{*subtask})
		glog.Info("subtask should be pushed back: ", subtask)
		return &errordef.DummyError{}
	}

	err = scheduler.AfterDispatch(subtask.SubtaskId)
	if err != nil {
		glog.Warning("AfterDispatch failed: ", subtask.SubtaskId, ",", err)
	}

	glog.Info("succeeded to dispatch subtask: ", subtask.SubtaskId, subtask.TaskId)
	return nil
}