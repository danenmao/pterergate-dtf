package dtfdef

import (
	"pterergate-dtf/dtf/extconfig"
)

// 服务角色，指定不同的服务类型
type ServiceRole uint32

const (
	ServiceRole_Manager   ServiceRole = 1
	ServiceRole_Generator ServiceRole = 2
	ServiceRole_Scheduler ServiceRole = 3
	ServiceRole_Executor  ServiceRole = 4
	ServiceRole_Iterator  ServiceRole = 5
)

// 服务配置
type ServiceConfig struct {
	MySQLServer     extconfig.MySQLAddress
	RedisServer     extconfig.RedisAddress
	ExecutorService extconfig.RPCServiceAddress
	IteratorService extconfig.RPCServiceAddress
}
