package routine

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
)

func Test_ExecRoutineWithInterval_NotifyToExit(t *testing.T) {
	counter := 0
	exitctrl.Register()
	go func() {
		time.Sleep(time.Second)
		exitctrl.NotifyToExit()
	}()

	ExecRoutineWithInterval("test",
		func() {
			counter += 1
		},
		100*time.Millisecond,
	)

	Convey("test to execute the working routine", t, func() {
		Convey("counter should be greater than 9", func() {
			So(counter, ShouldBeGreaterThan, 9)
		})
	})
}

func Test_StartWorkingRoutine_NotifyToExit(t *testing.T) {
	counter := 0
	exitctrl.Register()
	go func() {
		time.Sleep(time.Second)
		exitctrl.NotifyToExit()
	}()

	StartWorkingRoutine([]WorkingRoutine{
		{
			RoutineFn: func() {
				counter += 1
			},
			RoutineCount: 1,
			Interval:     100 * time.Millisecond,
		},
	})

	time.Sleep(time.Second)

	Convey("test to execute the working routine", t, func() {
		Convey("counter should be greater than 9", func() {
			So(counter, ShouldBeGreaterThan, 5)
		})
	})
}
