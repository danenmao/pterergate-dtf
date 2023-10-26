package schedulerflow

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow/executorconnector"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow/quotagroup"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/subtaskqueue"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskloader"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

func AddTaskToScheduler(
	taskId taskmodel.TaskIdType,
	groupName string,
	taskType uint32,
	priority uint32,
) error {
	return quotagroup.GetQuotaGroupMgr().AddTask(groupName, taskId, taskType, priority)
}

// if no task, retTaskId is 0, subtasks is empty
func ScheduleSubtasks(
	retTaskId *taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
) error {
	err := quotagroup.GetQuotaGroupMgr().Select(retTaskId, subtasks)
	if err != nil {
		glog.Warning("failed to select subtasks: ", err)
		return err
	}

	return nil
}

// execute subtasks belonging to the same task
func ExecSubtasks(
	taskId taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
	toPushbackSubtask *[]taskmodel.SubtaskBody,
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
	var scheduler taskmodel.ITaskSchedulerCallback = nil
	err = GetTaskSchedulerCallback(taskType, &scheduler)
	if err != nil {
		glog.Warning("failed to get task scheduler: ", taskType, ", ", err.Error())
		return err
	}

	glog.Info("get task type plugin scheduler: ", taskType, ", ", scheduler)

	// to dispatch subtasks
	doneSubtaskList := []taskmodel.SubtaskBody{}
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

	// to execute subtasks
	err = executorconnector.ExecSubtasks(taskId, subtasks)
	if err != nil {
		glog.Error("failed to execute subtasks: ", taskId, ", ", err.Error())
	}

	return nil
}

func GetTaskSchedulerCallback(taskType uint32, callback *taskmodel.ITaskSchedulerCallback) error {
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

	*callback = taskBody.SchedulerCallback
	glog.Info("succeeded to get task scheduler callback: ", taskType)
	return nil
}

// 对子任务执行调度操作
func DispatchSubtask(
	taskType uint32,
	callback taskmodel.ITaskSchedulerCallback,
	subtask *taskmodel.SubtaskBody,
) error {
	glog.Info("dipatch subtask: ", subtask)

	// invoke the dispatch method
	toDipatch, err := callback.BeforeDispatch(subtask.SubtaskId, subtask)
	if err != nil {
		glog.Info("dispatch subtask error,  pushed back: ", subtask, ", ", err.Error())
		return err
	}

	// don't dispatch it now, push it back
	if !toDipatch {
		subtaskqueue.PushSubtaskBack(subtask.TaskId, &[]taskmodel.SubtaskBody{*subtask})
		glog.Info("subtask should be pushed back: ", subtask)
		return &errordef.DummyError{}
	}

	err = callback.AfterDispatch(subtask.SubtaskId)
	if err != nil {
		glog.Warning("AfterDispatch failed: ", subtask.SubtaskId, ",", err)
	}

	glog.Info("succeeded to dispatch subtask: ", subtask.SubtaskId, subtask.TaskId)
	return nil
}
