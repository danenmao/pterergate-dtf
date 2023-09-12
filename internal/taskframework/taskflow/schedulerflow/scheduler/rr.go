package scheduler

import "pterergate-dtf/dtf/taskmodel"

// RR调度算法
type RR struct {
	QueueKeyName string
}

// RR调度算法
func (scheduler *RR) Schedule(
	retTaskId *taskmodel.TaskIdType,
	noTask *bool,
) error {
	return ScheduleQueue(scheduler.QueueKeyName, retTaskId, noTask)
}
