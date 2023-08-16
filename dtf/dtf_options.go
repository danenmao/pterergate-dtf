package dtf

import (
	"pterergate-dtf/dtf/config"
)

func WithMySQL(mysql *config.MySQLAddress) ServiceOptions {
	return func(config *ServiceConfig) {
		config.MySQLServer = *mysql
	}
}

func WithRedis(redis *config.RedisAddress) ServiceOptions {
	return func(config *ServiceConfig) {
		config.RedisServer = *redis
	}
}

func WithExecutor(executor *config.RPCServiceAddress) ServiceOptions {
	return func(config *ServiceConfig) {
		config.ExecutorService = *executor
	}
}

func WithIterator(iterator *config.RPCServiceAddress) ServiceOptions {
	return func(config *ServiceConfig) {
		config.IteratorService = *iterator
	}
}
