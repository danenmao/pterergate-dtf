package servicectrl

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/dtfdef"
	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
)

// start the specified service
func StartService(role dtfdef.ServiceRole, cfg *dtfdef.ServiceConfig) error {

	// search service role start fn
	starter, found := gs_ServiceRoleStarter[role]
	if !found {
		glog.Warning("unknown service role: ", role)
		return errordef.ErrInvalidParameter
	}

	// to process the exit signal
	exitctrl.RegisterWithDuration(cfg.PrestopDuration)

	// invoke the start fn
	starter(cfg)

	return nil
}

// notify to stop all routines
func NotifyStop() error {
	exitctrl.NotifyToExit()
	return nil
}

// wait for all routines to exit
func Join() error {
	exitctrl.Join()
	return nil
}
