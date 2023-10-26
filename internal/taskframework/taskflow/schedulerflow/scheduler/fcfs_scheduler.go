package scheduler

import (
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

// FCFS调度算法
type FCFSScheduler struct {
	QueueKeyName string // 队列的Key名
}

// FCFS调度算法
func (scheduler *FCFSScheduler) Schedule(
	retTaskId *taskmodel.TaskIdType,
	noTask *bool,
) error {
	return ScheduleTaskInQueue(scheduler.QueueKeyName, retTaskId, noTask)
}
