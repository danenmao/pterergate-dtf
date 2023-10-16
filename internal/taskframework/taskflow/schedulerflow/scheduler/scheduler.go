package scheduler

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

type IScheduleQueueImpl interface {
	Schedule(taskId *taskmodel.TaskIdType, noTask *bool) error
}
