package signalctrl

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

var (
	// 退出标记。收到退出信号时设置此标记
	s_NotifyToExitFlag = false

	// 真实的退出标记
	s_JustExitFlag = false
)

// 用于全局性通知退出的context
var SignalContext context.Context = nil
var s_SignalCancelFn context.CancelFunc = nil

// 注册处理退出信号
func RegisterSignal() error {
	SignalContext, s_SignalCancelFn = context.WithCancel(context.Background())
	go listenSignal()
	return nil
}

// 判断是否需要退出
func CheckIfNeedToExit() bool {
	return s_NotifyToExitFlag
}

// 等待退出
func WaitPreStop() {
	for {
		if s_JustExitFlag {
			return
		}

		select {
		case <-SignalContext.Done():
			return
		default:
			time.Sleep(time.Second)
		}
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
			doExit()
			glog.Warning("exited")

		default:
			glog.Warning("unknown signal: ", s)
		}

	}
}

// 执行退出操作
func doExit() {

	// 设置退出标记
	s_NotifyToExitFlag = true

	// 通知服务退出
	s_SignalCancelFn()

	// 等待一段时间
	// 配合一些服务发现机制的更新间隔, 在接收到退出信号后，等待若干秒再退出
	time.Sleep(time.Duration(PreStopWaitInterval) * time.Second)

	// 返回，执行退出操作
	s_JustExitFlag = true
}
