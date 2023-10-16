package servicectrl

import (
	"time"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/services/collector"
)

func StartCollector(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    collector.CompleteSubtaskRoutine,
			RoutineCount: config.EnvCompleteSubtaskConcurrencyLimit,
			Interval:     time.Millisecond * time.Duration(config.EnvCompleteSubtaskInterval),
		},
	})

	// register collector handler
	cfg.CollectorHandlerRegister(collector.CollectorRequestHandler)

	return nil
}
