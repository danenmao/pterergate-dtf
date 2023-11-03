package serversupport

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

func TestMain(m *testing.M) {
	fmt.Println("setup...")
	Setup()

	retCode := m.Run()

	fmt.Println("teardown...")
	Teardown()
	os.Exit(retCode)
}

func Test_CollectorGetInvoker_Success(t *testing.T) {
	requestFlag := false
	svr := NewCollectorServer(
		func([]taskmodel.SubtaskResult) error {
			requestFlag = true
			return nil
		})

	go func() {
		svr.Serve(8090)
	}()
	time.Sleep(10 * time.Millisecond)

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
