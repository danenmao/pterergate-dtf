package serverhelper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	fmt.Println("setup...")
	Setup()

	retCode := m.Run()

	fmt.Println("teardown...")
	Teardown()
	os.Exit(retCode)
}

func Test_Post_Success(t *testing.T) {
	requestFlag := false
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			requestFlag = true
			return "", nil
		}})

	go func() {
		svr.StartServer()
	}()

	invoker := NewSimpleInvoker()
	err := invoker.Post("http://localhost:8090/test", "test", "test")
	svr.Shutdown()

	Convey("post a request successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeTrue)
		})
	})
}

func Test_Post_ResponseError(t *testing.T) {
	requestFlag := false
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			requestFlag = true
			return "testerror", errors.New("test error")
		}})

	go func() {
		svr.StartServer()
	}()

	invoker := NewSimpleInvoker()
	err := invoker.Post("http://localhost:8090/test", "test", "test")
	svr.Shutdown()

	Convey("a request return a error", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
			So(requestFlag, ShouldBeTrue)
		})
	})
}

func Test_Post_InvalidUrl(t *testing.T) {
	invoker := NewSimpleInvoker()
	err := invoker.Post("http://localhost:8090/test", "test", "test")

	Convey("failed to post a request to an invalid url", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_Post_InvalidCommonResponse(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.POST(
		"/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, "test response")
		},
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8090),
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Print("failed to run the gin server: ", err.Error())
		}
	}()

	invoker := NewSimpleInvoker()
	err := invoker.Post("http://localhost:8090/test", "test", "test")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	server.Shutdown(ctx)

	Convey("invalid CommonResponse", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
