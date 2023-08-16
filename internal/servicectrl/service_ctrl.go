package servicectrl

import (
	"github.com/golang/glog"

	"pterergate-dtf/dtf/dtfdef"
	"pterergate-dtf/dtf/errordef"
)

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, cfg *dtfdef.ServiceConfig) error {

	startFn, found := gs_ServiceRoleAction[role]
	if !found {
		glog.Warning("unknown service role: ", role)
		return errordef.ErrInvalidParameter
	}

	startFn(cfg)

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
