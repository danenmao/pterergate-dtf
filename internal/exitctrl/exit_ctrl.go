package exitctrl

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

// preStop间隔，收到退出信号时的等待间隔
var PreStopWaitInterval = 10

// 用于全局性通知退出的context
var SignalContext context.Context = nil

type ExitSignalController struct {
	NotifyToExitFlag bool // 通知退出的标记, 收到退出信号时设置此标记
	JustExitFlag     bool // 立即退出的标记
	CancelFn         context.CancelFunc
}

var gs_ExitController = ExitSignalController{
	NotifyToExitFlag: false,
	JustExitFlag:     false,
	CancelFn:         nil,
}

// 注册处理退出信号
func Register() error {
	SignalContext, gs_ExitController.CancelFn = context.WithCancel(context.Background())
	go listenSignal()
	return nil
}

// 通知退出
func NotifyToExit() {
	if gs_ExitController.CancelFn == nil {
		panic("no valid cancel fn")
	}

	gs_ExitController.NotifyToExitFlag = true
	gs_ExitController.CancelFn()
}

// 判断是否需要退出
func CheckIfNeedToExit() bool {
	return gs_ExitController.NotifyToExitFlag
}

// 等待退出通知
func WaitForNotify(interval time.Duration) bool {
	select {
	case <-SignalContext.Done():
		return true
	default:
		time.Sleep(interval)
		return false
	}
}

// 等待退出
func Prestop() {
	for {
		if gs_ExitController.JustExitFlag {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// 注册并监听注册的退出信号
func listenSignal() {
	// 注册监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// 等待信号
	for s := range c {
		glog.Warning("get a signal: ", s)

		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Warning("to exit")
			clean()
			glog.Warning("exited")

		default:
			glog.Warning("unknown signal: ", s)
		}
	}
}

// 执行退出操作
func clean() {
	// 设置退出标记
	gs_ExitController.NotifyToExitFlag = true

	// 通知服务退出
	gs_ExitController.CancelFn()

	// 等待一段时间
	// 配合一些服务发现机制的更新间隔, 在接收到退出信号后，等待若干秒再退出
	time.Sleep(time.Duration(PreStopWaitInterval) * time.Second)

	// 返回，执行退出操作
	gs_ExitController.JustExitFlag = true
}
