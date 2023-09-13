package dtf

import (
	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/dtf/extconfig"
	"pterergate-dtf/dtf/taskmodel"
)

// 用于设置服务配置
type ServiceOption func(config *dtfdef.ServiceConfig)

func WithMySQL(mysql *extconfig.MySQLAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.MySQLServer = *mysql
	}
}

func WithRedis(redis *extconfig.RedisAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.RedisServer = *redis
	}
}

func WithMongoDB(mongo *extconfig.MongoAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.MongoServer = *mongo
	}
}

func WithExecutor(executor taskmodel.ExecutorInvoker) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.ExecutorService = executor
	}
}

func WithCollector(collector taskmodel.CollectorInvoker) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.CollectorService = collector
	}
}

func WithRegisterExecutorHandler(register taskmodel.RegisterExecutorRequestHandler) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.ExecutorHandlerRegister = register
	}
}

func WithRegisterCollectorHandler(register taskmodel.RegisterCollectorRequestHandler) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.CollectorHandlerRegister = register
	}
}
