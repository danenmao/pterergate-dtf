package servicectrl

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/services/scheduler"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow/executorconnector"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow/resourcegroup"
)

func StartScheduler(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	executorconnector.ExecutorService = cfg.ExecutorService
	resourcegroup.GetResourceGroupMgr().Init()

	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    scheduler.ScheduleTaskRoutine,
			RoutineCount: config.EnvScheduleTaskConcurrencyLimit,
			Interval:     time.Duration(config.EnvScheduleTaskInterval) * time.Millisecond,
		},
		{
			RoutineFn:    scheduler.MonitorSubtaskComplete,
			RoutineCount: config.EnvMonitorSubtaskCompleteConcurrencyLimit,
			Interval:     time.Millisecond * time.Duration(config.EnvMonitorSubtaskCompleteInterval),
		},
		{
			RoutineFn:    scheduler.MonitorTimeoutSubtask,
			RoutineCount: config.EnvMonitorSubtaskTimeoutConcurrencyLimit,
			Interval:     time.Second * time.Duration(config.EnvMonitorSubtaskTimeoutInterval),
		},
		{
			RoutineFn:    scheduler.MonitorRunningTaskToComplete,
			RoutineCount: config.EnvMonitorTaskCompleteConcurrencyLimit,
			Interval:     time.Second * time.Duration(config.EnvMonitorTaskCompleteInterval),
		},
	})

	return nil
}
