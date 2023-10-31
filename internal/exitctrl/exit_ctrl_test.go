package exitctrl

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Register_Normal(t *testing.T) {
	Register()
	NotifyToExit()

	Convey("need to exit", t, func() {
		Convey("should be true", func() {
			So(IfNeedToExit(), ShouldBeTrue)
		})
	})

	WaitForSignal(500 * time.Millisecond)
}

func Test_IfNeedToExit_False(t *testing.T) {
	Register()

	Convey("IfNeedToExit returns false", t, func() {
		Convey("should be false", func() {
			So(IfNeedToExit(), ShouldBeFalse)
		})
	})

	NotifyToExit()
	WaitForSignal(500 * time.Millisecond)
}

func Test_WaitForSignal_True(t *testing.T) {
	Register()
	NotifyToExit()

	Convey("WaitForSignal returns true", t, func() {
		Convey("should be true", func() {
			So(WaitForSignal(300*time.Millisecond), ShouldBeTrue)
		})
	})
}

func Test_WaitForSignal_False(t *testing.T) {
	Register()

	Convey("WaitForSignal returns false", t, func() {
		Convey("should be false", func() {
			So(WaitForSignal(100*time.Millisecond), ShouldBeFalse)
		})
	})

	NotifyToExit()
	WaitForSignal(500 * time.Millisecond)
}

func Test_clean_True(t *testing.T) {
	Register()
	NotifyToExit()
	WaitForSignal(500 * time.Millisecond)

	Convey("test clean function", t, func() {
		Convey("should be true", func() {
			So(gs_Controller.JustExitFlag, ShouldBeTrue)
		})
	})
}

func Test_Prestop_Normal(t *testing.T) {
	RegisterWithDuration(400 * time.Millisecond)
	NotifyToExit()
	Prestop()

	Convey("test clean function", t, func() {
		Convey("should be true", func() {
			So(gs_Controller.JustExitFlag, ShouldBeTrue)
		})
	})
}
