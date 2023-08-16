package servicectrl

import (
	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
)

func StartManager(cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	return nil
}
