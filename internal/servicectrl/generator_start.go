package servicectrl

import (
	"time"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/idtool"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/services/generator"
)

func StartGenerator(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	// init dependencies
	idtool.Init(config.TaskIdKey)

	// start service working routines
	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    generator.StartGenerateTaskRoutine,
			RoutineCount: config.EnvGenerateTaskConcurrencyLimit,
			Interval:     time.Second * time.Duration(config.EnvGenerateTaskCheckInterval),
		},
		{
			RoutineFn:    generator.MonitorTaskGenerationRoutine,
			RoutineCount: config.EnvMonitorTaskGenerationConcurrencyLimit,
			Interval:     time.Second * time.Duration(config.EnvMonitorTaskGenerationInterval),
		},
	})

	return nil
}
