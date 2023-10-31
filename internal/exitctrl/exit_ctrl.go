package exitctrl

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

// prestop interval
var gs_PreStop = 10 * time.Second

// the context to be notified to exit
var SignalCtx context.Context = nil
var gs_ExitChan chan os.Signal = nil

type ExitSignalController struct {
	NotifyToExitFlag bool // 通知退出的标记, 收到退出信号时设置此标记
	JustExitFlag     bool // 立即退出的标记
	CancelFn         context.CancelFunc
}

var gs_Controller = ExitSignalController{
	NotifyToExitFlag: false,
	JustExitFlag:     false,
	CancelFn:         nil,
}

// register to process the exit signal
func Register() error {
	RegisterWithDuration(0)
	return nil
}

func RegisterWithDuration(duration time.Duration) error {
	// reset the state
	gs_PreStop = duration
	gs_Controller.NotifyToExitFlag = false
	gs_Controller.JustExitFlag = false
	SignalCtx, gs_Controller.CancelFn = context.WithCancel(context.Background())

	// register the signal
	gs_ExitChan = make(chan os.Signal, 1)
	signal.Notify(
		gs_ExitChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// listen
	go listenToSignal()
	return nil
}

// notify the routines  to exit
func NotifyToExit() {
	gs_ExitChan <- os.Interrupt
	gs_Controller.NotifyToExitFlag = true
}

// check if the caller need to exit
func IfNeedToExit() bool {
	return gs_Controller.NotifyToExitFlag
}

// wait for the exit signal
func WaitForSignal(interval time.Duration) bool {
	now := time.Now()
	for {
		select {
		case <-SignalCtx.Done():
			return true

		default:
			if time.Since(now) >= interval {
				return false
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}

// prestop function
func Prestop() {
	for {
		if gs_Controller.JustExitFlag {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// register the signal and listen
func listenToSignal() {
	// wait for the signal
	for s := range gs_ExitChan {
		glog.Warning("get a signal: ", s)

		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Warning("to exit")
			clean()
			glog.Warning("exited")
			return

		default:
			glog.Warning("unknown signal: ", s)
		}
	}
}

// perform the clean operation
func clean() {
	// set the notify flag
	gs_Controller.NotifyToExitFlag = true
	gs_Controller.CancelFn()

	// prestop
	time.Sleep(gs_PreStop)
	gs_Controller.JustExitFlag = true
}
