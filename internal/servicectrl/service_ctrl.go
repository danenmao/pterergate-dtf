package servicectrl

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/signalctrl"
)

// 启动指定的服务
func StartService(role dtfdef.ServiceRole, cfg *dtfdef.ServiceConfig) error {

	// search service role start fn
	startFn, found := gs_ServiceRoleAction[role]
	if !found {
		glog.Warning("unknown service role: ", role)
		return errordef.ErrInvalidParameter
	}

	// to process exit signal
	signalctrl.RegisterSignal()

	// invoke the start fn
	startFn(cfg)

	return nil
}

// 通知停止服务
func NotifyStop() error {
	signalctrl.NotifyToExit()
	return nil
}

// 等待服务停止
func Join() error {
	signalctrl.WaitPreStop()
	return nil
}
