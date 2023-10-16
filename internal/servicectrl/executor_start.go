package servicectrl

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/mysqltool"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/services/executor"
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
