package scheduler

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/subtaskqueue"
)

const (
	// subtask timeout, second
	s_timeout_second = 900
)

// go_schedule_subtask
func ScheduleTaskRoutine() {

	// get the task and subtasks to schedule
	var taskId taskmodel.TaskIdType
	var subtasks = []taskmodel.SubtaskBody{}
	err := schedulerflow.ScheduleSubtasks(&taskId, &subtasks)
	if err != nil {
		glog.Warning("failed to get task to schedule: ", err)
		return
	}

	// no task, return
	if len(subtasks) <= 0 {
		return
	}

	glog.Info("to exec subtasks of task: ", taskId, ", ", len(subtasks))

	// to execute subtasks
	toPushbackSubtask := []taskmodel.SubtaskBody{}
	err = schedulerflow.ExecSubtasks(taskId, &subtasks, &toPushbackSubtask)
	if err != nil {
		glog.Warning("failed to exec subtask: ", taskId, ",", err)

		// if failed, push back all subtasks
		subtaskqueue.PushSubtaskBack(taskId, &subtasks)
		return
	}

	// push back all subtasks
	if len(toPushbackSubtask) > 0 {
		glog.Info("push subtasks back to generation queue: ", taskId, ", ", len(toPushbackSubtask))
		subtaskqueue.PushSubtaskBack(taskId, &toPushbackSubtask)
	}

	glog.Info("succeeded to schedule subtasks: ", taskId)
}
