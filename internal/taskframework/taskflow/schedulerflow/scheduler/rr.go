package scheduler

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

// RR调度算法
type RRScheduler struct {
	QueueKeyName string
}

// RR调度算法
func (scheduler *RRScheduler) Schedule(
	retTaskId *taskmodel.TaskIdType,
	noTask *bool,
) error {
	return ScheduleQueue(scheduler.QueueKeyName, retTaskId, noTask)
}
