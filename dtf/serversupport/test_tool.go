package serversupport

import (
	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
	"github.com/danenmao/pterergate-dtf/internal/serverhelper"
)

func Setup() {
	exitctrl.Register()
	serverhelper.Setup()
}

func Teardown() {
	exitctrl.NotifyToExit()
	serverhelper.Teardown()
	exitctrl.Join()
}
