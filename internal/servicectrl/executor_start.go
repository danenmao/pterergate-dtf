package servicectrl

import (
	"time"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/routine"
	"pterergate-dtf/internal/services/executor"
)

func StartExecutor(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	executor.CollectorInvoker = cfg.CollectorService
	executor.GetExecutorService().Init()
	routine.StartWorkingRoutine([]routine.WorkingRoutine{
		{
			RoutineFn:    executor.ReportRoutine,
			RoutineCount: 1,
			Interval:     time.Second,
		},
	})

	// register request handler
	cfg.ExecutorHandlerRegister(executor.ExecutorRequestHandler)

	return nil
}
