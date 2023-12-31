package servicectrl

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/idtool"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/services/taskmgmt"
)

func StartManager(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	// init dependencies
	idtool.Init(config.TaskIdKey)

	// start service working routines
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
