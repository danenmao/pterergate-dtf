package scheduler

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

type IScheduleQueue interface {
	Schedule(taskId *taskmodel.TaskIdType, noTask *bool) error
}
