package servicectrl

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/idtool"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/services/generator"
)

func StartGenerator(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	// init dependencies
	idtool.Init(config.SubtaskIdKey)

	// start service working routines
	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    generator.StartTaskGenerationRoutine,
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
