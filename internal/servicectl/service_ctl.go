package servicectl

import (
	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
)

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, cfg *dtfdef.ServiceConfig) error {

	config.DefaultMySQL = cfg.MySQLServer
	mysqltool.ConnectToDefaultMySQL()

	config.DefaultRedisServer = cfg.RedisServer
	redistool.ConnectToDefaultRedis()

	return nil
}

// 通知停止服务
func NotifyStop() error {
	return nil
}

// 等待服务停止
func Join() error {
	return nil
}
