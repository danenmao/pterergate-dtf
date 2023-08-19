package dtf

import (
	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/dtf/extconfig"
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

func WithExecutor(executor *extconfig.RPCServiceAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.ExecutorService = *executor
	}
}

func WithIterator(iterator *extconfig.RPCServiceAddress) ServiceOption {
	return func(config *dtfdef.ServiceConfig) {
		config.IteratorService = *iterator
	}
}
