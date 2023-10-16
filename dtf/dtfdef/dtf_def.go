package dtfdef

import (
	"github.com/danenmao/pterergate-dtf/dtf/extconfig"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

// 服务角色，指定不同的服务类型
type ServiceRole uint32

const (
	ServiceRole_Manager   ServiceRole = 1
	ServiceRole_Generator ServiceRole = 2
	ServiceRole_Scheduler ServiceRole = 3
	ServiceRole_Executor  ServiceRole = 4
	ServiceRole_Collector ServiceRole = 5
)

// 服务配置
type ServiceConfig struct {
	MySQLServer              extconfig.MySQLAddress
	RedisServer              extconfig.RedisAddress
	MongoServer              extconfig.MongoAddress
	ExecutorService          taskmodel.ExecutorInvoker
	CollectorService         taskmodel.CollectorInvoker
	ExecutorHandlerRegister  taskmodel.RegisterExecutorRequestHandler
	CollectorHandlerRegister taskmodel.RegisterCollectorRequestHandler
}
