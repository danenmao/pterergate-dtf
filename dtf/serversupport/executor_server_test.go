package serversupport

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

func Test_ExecutorGetRegister_Success(t *testing.T) {
	requestFlag := false
	svr := NewExecutorServer(nil)
	svr.GetRegister()(func([]taskmodel.SubtaskBody) error {
		requestFlag = true
		return nil
	})

	go func() {
		svr.Serve(8090)
	}()
	time.Sleep(10 * time.Millisecond)

	invoker := NewExecutorInvoker("localhost", 8090, "test-invoker")
	err := invoker.GetInvoker()([]taskmodel.SubtaskBody{})
	svr.Shutdown()

	Convey("post a request successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeTrue)
		})
	})
}
