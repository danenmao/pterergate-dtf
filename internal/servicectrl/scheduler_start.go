package servicectrl

import (
	"time"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/services/scheduler"
	"pterergate-dtf/internal/taskframework/taskflow/schedulerflow/executorconnector"
	"pterergate-dtf/internal/taskframework/taskflow/schedulerflow/resourcegroup"
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
			RoutineFn:    scheduler.MonitorSubtaskTimeout,
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
