package dtf

import (
	"pterergate-dtf/dtf/config"
	"pterergate-dtf/dtf/dtfdef"
)

// 用于设置服务配置
type ServiceOption func(config *dtfdef.ServiceConfig)

func WithMySQL(mysql *config.MySQLAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.MySQLServer = *mysql
	}
}

func WithRedis(redis *config.RedisAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.RedisServer = *redis
	}
}

func WithExecutor(executor *config.RPCServiceAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.ExecutorService = *executor
	}
}

func WithIterator(iterator *config.RPCServiceAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.IteratorService = *iterator
	}
}
