package idtool

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

func TestMain(m *testing.M) {
	fmt.Println("setup...")
	redistool.Setup()

	retCode := m.Run()

	fmt.Println("teardown...")
	redistool.Teardown()
	os.Exit(retCode)
}

func Test_Init_NoKeyname(t *testing.T) {
	keeper := IdKeeper{}
	err := keeper.Init("")
	Convey("get the id", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_GetId_NoKeyname(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	id, err := keeper.GetId("")
	exitctrl.NotifyToExit()

	Convey("get the id", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
			So(id, ShouldEqual, 0)
		})
	})

	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_GetId_DifferentKeyname(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	id, err := keeper.GetId("different")
	exitctrl.NotifyToExit()

	Convey("get the id", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
			So(id, ShouldEqual, 0)
		})
	})

	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_GetId_Initial(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)

	id, err := keeper.GetId(keyName)
	Convey("get the id", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 1)
			So(keeper.Count, ShouldEqual, ReallocStep-1)
			So(keeper.Start, ShouldEqual, 2)
			So(keeper.End, ShouldEqual, ReallocStep+1)
			So(keeper.FormerEnd, ShouldBeZeroValue)
			So(keeper.NewStart, ShouldEqual, 1)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_GetId_Fail(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetErr(errors.New("test error"))

	id, err := keeper.GetId(keyName)
	Convey("get the id", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
			So(id, ShouldEqual, 0)
			So(keeper.Count, ShouldEqual, 0)
			So(keeper.Start, ShouldEqual, 0)
			So(keeper.End, ShouldEqual, 0)
			So(keeper.FormerEnd, ShouldBeZeroValue)
			So(keeper.NewStart, ShouldEqual, 0)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_realloc_NeedNotToAlloc(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)
	keeper.GetId(keyName)

	keeper.realloc()

	Convey("get the id", t, func() {
		Convey("should be nil", func() {
			So(keeper.Count, ShouldEqual, ReallocStep-1)
			So(keeper.Start, ShouldEqual, 2)
			So(keeper.End, ShouldEqual, ReallocStep+1)
			So(keeper.FormerEnd, ShouldBeZeroValue)
			So(keeper.NewStart, ShouldEqual, 1)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_realloc_AdjacentDualRange(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)

	for i := 0; i < ReallocStep-ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	expInt = redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(2 * ReallocStep)
	keeper.realloc()

	id, err := keeper.GetId(keyName)
	Convey("realloc an id range", t, func() {
		Convey("alloc an adjacent range", func() {
			So(err, ShouldBeNil)
			So(id, ShouldEqual, ReallocStep-ReallocThreshold+1)
			So(keeper.Count, ShouldEqual, ReallocStep+ReallocThreshold-1)
			So(keeper.Start, ShouldEqual, ReallocStep-ReallocThreshold+2)
			So(keeper.End, ShouldEqual, 2*ReallocStep+1)
			So(keeper.FormerEnd, ShouldEqual, ReallocStep+1)
			So(keeper.NewStart, ShouldEqual, ReallocStep+1)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_realloc_MoveToAdjacentRange(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)

	for i := 0; i < ReallocStep-ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	expInt = redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(2 * ReallocStep)
	keeper.realloc()

	for i := 0; i < ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	id, err := keeper.GetId(keyName)
	Convey("move to an adjacent range", t, func() {
		Convey("move to an adjacent range", func() {
			So(err, ShouldBeNil)
			So(id, ShouldEqual, ReallocStep+1)
			So(keeper.Count, ShouldEqual, ReallocStep-1)
			So(keeper.Start, ShouldEqual, ReallocStep+2)
			So(keeper.End, ShouldEqual, 2*ReallocStep+1)
			So(keeper.FormerEnd, ShouldEqual, 0)
			So(keeper.NewStart, ShouldEqual, 0)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_realloc_NonadjacentDualRange(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)

	for i := 0; i < ReallocStep-ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	expInt = redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(5 * ReallocStep)
	keeper.realloc()

	id, err := keeper.GetId(keyName)
	Convey("realloc a nonadjacent range", t, func() {
		Convey("alloc a nonadjacent range", func() {
			So(err, ShouldBeNil)
			So(id, ShouldEqual, ReallocStep-ReallocThreshold+1)
			So(keeper.Count, ShouldEqual, ReallocStep+ReallocThreshold-1)
			So(keeper.Start, ShouldEqual, ReallocStep-ReallocThreshold+2)
			So(keeper.End, ShouldEqual, 5*ReallocStep+1)
			So(keeper.FormerEnd, ShouldEqual, ReallocStep+1)
			So(keeper.NewStart, ShouldEqual, 4*ReallocStep+1)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}

func Test_realloc_MoveToNonadjacentRange(t *testing.T) {
	exitctrl.Register()

	keyName := "test"
	keeper := IdKeeper{}
	keeper.Init(keyName)

	expInt := redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(ReallocStep)

	for i := 0; i < ReallocStep-ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	expInt = redistool.ClientMock.ExpectIncrBy(keyName, ReallocStep)
	expInt.SetVal(5 * ReallocStep)
	keeper.realloc()

	for i := 0; i < ReallocThreshold; i++ {
		keeper.GetId(keyName)
	}

	id, err := keeper.GetId(keyName)
	Convey("move to a nonadjacent range", t, func() {
		Convey("move to a nonadjacent range", func() {
			So(err, ShouldBeNil)
			So(id, ShouldEqual, 4*ReallocStep+1)
			So(keeper.Count, ShouldEqual, ReallocStep-1)
			So(keeper.Start, ShouldEqual, 4*ReallocStep+2)
			So(keeper.End, ShouldEqual, 5*ReallocStep+1)
			So(keeper.FormerEnd, ShouldEqual, 0)
			So(keeper.NewStart, ShouldEqual, 0)
		})
	})

	exitctrl.NotifyToExit()
	exitctrl.WaitForSignal(200 * time.Millisecond)
}
