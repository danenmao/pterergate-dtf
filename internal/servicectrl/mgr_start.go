package servicectrl

import (
	"time"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/idtool"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/services/taskmgmt"
)

func StartManager(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	idtool.Init(config.TaskIdKey)
	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    taskmgmt.MonitorTaskTableRoutine,
			RoutineCount: config.EnvMonitorTaskTblCountLimit,
			Interval:     time.Duration(config.EnvMonitorTaskTblInterval) * time.Second,
		},
		{
			RoutineFn:    taskmgmt.MonitorTaskTimeout,
			RoutineCount: config.EnvMonitorTaskTimeoutCountLimit,
			Interval:     time.Duration(config.EnvMonitorTaskTimeoutInterval) * time.Second,
		},
		{
			RoutineFn:    taskmgmt.MonitorCompletedTask,
			RoutineCount: config.EnvMonitorTaskCompletedCountLimit,
			Interval:     time.Duration(config.EnvMonitorTaskCompletedInterval) * time.Second,
		},
	})

	return nil
}
