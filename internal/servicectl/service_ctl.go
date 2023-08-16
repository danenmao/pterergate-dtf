package servicectl

import (
	"pterergate-dtf/dtf/dtfdef"
)

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, cft *dtfdef.ServiceConfig) error {
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
