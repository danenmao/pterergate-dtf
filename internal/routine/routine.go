package routine

import (
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
	"github.com/danenmao/pterergate-dtf/internal/misc"
)

// 例程类型
type RoutineFn func()

// 工作例程结构
type WorkingRoutine struct {
	RoutineFn    RoutineFn     // 工作例程函数
	RoutineCount uint          // 例程数量
	Interval     time.Duration // 工作例程的执行间隔
}

// 启动所有的工作例程
func StartWorkingRoutine(workers []WorkingRoutine) error {
	// 按数量创建工作例程
	for _, worker := range workers {
		name := misc.GetFunctionName(worker.RoutineFn)
		for i := 0; i < int(worker.RoutineCount); i++ {
			go WorkingRoutineWrapper(name, worker.RoutineFn, worker.Interval)()
		}
	}

	return nil
}

// 工作例程包装函数
func WorkingRoutineWrapper(name string, fn RoutineFn, interval time.Duration) RoutineFn {
	return func() {
		ExecRoutineWithInterval(name, fn, interval)
	}
}

// 定期执行工作例程
func ExecRoutineWithInterval(
	name string,
	routine RoutineFn,
	interval time.Duration,
) {
	glog.Info("begin to ", name)

	// 定期执行检查
	for {
		// 检查并等待退出信号
		if exitctrl.WaitForSignal(interval) {
			glog.Info("got to exit signal")
			break
		}

		// 执行例程
		routine()
	}

	glog.Info("leave ", name)
}
