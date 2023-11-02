package serversupport

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

func Test_CollectorGetRegister_Success(t *testing.T) {
	requestFlag := false
	svr := NewCollectorServer(nil)
	svr.GetRegister()(func([]taskmodel.SubtaskResult) error {
		requestFlag = true
		return nil
	})

	go func() {
		svr.Serve(8090)
	}()

	invoker := NewCollectorInvoker("localhost", 8090, "test-invoker")
	err := invoker.GetInvoker()([]taskmodel.SubtaskResult{})
	svr.Shutdown()

	Convey("post a request successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeTrue)
		})
	})
}
