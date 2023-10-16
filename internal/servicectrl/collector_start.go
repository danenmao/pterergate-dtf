package servicectrl

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/services/collector"
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
