package serverhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/msgsigner"
)

func Test_NewSimpleServer_Success(t *testing.T) {
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			return "", nil
		}})

	Convey("new a server successfully", t, func() {
		Convey("should be nil", func() {
			So(svr, ShouldNotBeNil)
			So(svr.ServerPort == 8090, ShouldBeTrue)
			So(len(svr.handlerMap) == 1, ShouldBeTrue)
			So(svr.signer, ShouldNotBeNil)
		})
	})
}

func Test_NewSimpleServer_NilMap(t *testing.T) {
	svr := NewSimpleServer(8090, nil)

	Convey("new a server successfully", t, func() {
		Convey("should be nil", func() {
			So(svr, ShouldNotBeNil)
			So(svr.ServerPort == 8090, ShouldBeTrue)
			So(len(svr.handlerMap) == 0, ShouldBeTrue)
			So(svr.signer, ShouldNotBeNil)
		})
	})
}

func Test_Shutdown_NilServer(t *testing.T) {
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			return "", nil
		}})

	err := svr.Shutdown()
	Convey("new a server successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(svr, ShouldNotBeNil)
			So(svr.ServerPort == 8090, ShouldBeTrue)
			So(len(svr.handlerMap) == 1, ShouldBeTrue)
			So(svr.signer, ShouldNotBeNil)
		})
	})
}

func Test_Shutdown_Success(t *testing.T) {
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

	time.Sleep(10 * time.Millisecond)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	invoker := NewSimpleInvoker()
	_, err := invoker.Post("http://localhost:8090/test", "test", "test")

	Convey("The server shutdown successfully", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldNotBeNil)
			So(requestFlag, ShouldBeFalse)
		})
	})
}

func Test_Shutdown_Fail(t *testing.T) {
	requestFlag := false
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			requestFlag = true
			time.Sleep(30 * time.Millisecond)
			return "", nil
		}})

	go func() {
		svr.StartServer()
	}()
	time.Sleep(10 * time.Millisecond)

	invoker := NewSimpleInvoker()
	go func() {
		invoker.Post("http://localhost:8090/test", "test", "test")
	}()
	time.Sleep(10 * time.Millisecond)

	err := svr.ShutdownWithDuration(10 * time.Millisecond)
	Convey("The server failed to shutdown", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldNotBeNil)
			So(requestFlag, ShouldBeTrue)
		})
	})

	time.Sleep(50 * time.Millisecond)
}

func Test_authMiddleware_NoAuthorization(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader(string("test")))
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request without authentication", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_AuthorizationFailed)
		})
	})
}

func Test_authMiddleware_InvalidAuthorizationFormat(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader(string("test")))
	httpReq.Header.Set("Authorization", "TEST 1234567890")
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request with invalid authentication format", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_AuthorizationFailed)
		})
	})
}

func Test_authMiddleware_InvalidSign(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader(string("test")))
	httpReq.Header.Set("Authorization", "Bearer 1234567890")
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request with invalid signature", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_AuthorizationFailed)
		})
	})
}

func Test_authMiddleware_InvalidMsg(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader("test"))
	sign, _ := msgsigner.NewMsgSigner().Sign("", "", []string{""}, "test", time.Minute)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request with invalid msg", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_AuthorizationFailed)
		})
	})
}

func Test_handleRequest_DifferentHash(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	var request = CommonRequest{}
	request.Body = "body test data"
	request.Header.BodyHash = "1234567890"
	cr, _ := json.Marshal(&request)

	msg := CommonMessage{
		UserName: "test",
		BodyHash: request.Header.BodyHash,
	}
	msgPlain, _ := json.Marshal(msg)
	sign, _ := msgsigner.NewMsgSigner().Sign("", "", []string{""}, string(msgPlain), time.Minute)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader(string(cr)))
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request without authentication", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_AuthorizationFailed)
		})
	})
}

func Test_handleRequest_InvalidCommonRequest(t *testing.T) {
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
	time.Sleep(10 * time.Millisecond)

	msg := CommonMessage{
		UserName: "test",
		BodyHash: "1234567890",
	}
	msgPlain, _ := json.Marshal(msg)
	sign, _ := msgsigner.NewMsgSigner().Sign("", "", []string{""}, string(msgPlain), time.Minute)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader("test request data"))
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request without authentication", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeFalse)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_ParsingParam)
		})
	})
}

func Test_handleRequest_HandlerReturnError(t *testing.T) {
	requestFlag := false
	svr := NewSimpleServer(8090,
		map[string]RequestHandler{"/test": func(
			header RequestHeader, requestBody string) (responseBody string, err error) {
			requestFlag = true
			return "", errordef.ErrOperationFailed
		}})

	go func() {
		svr.StartServer()
	}()
	time.Sleep(10 * time.Millisecond)

	var request = CommonRequest{}
	request.Body = "body test data"
	request.Header.BodyHash = CalcMsgHash(request.Body)
	cr, _ := json.Marshal(&request)

	msg := CommonMessage{
		UserName: "test",
		BodyHash: request.Header.BodyHash,
	}
	msgPlain, _ := json.Marshal(msg)
	sign, _ := msgsigner.NewMsgSigner().Sign("", "", []string{""}, string(msgPlain), time.Minute)

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8090/test",
		strings.NewReader(string(cr)))
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))
	rsp, err := http.DefaultClient.Do(httpReq)

	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	commonResp := CommonResponse{}
	json.Unmarshal(body, &commonResp)
	svr.ShutdownWithDuration(100 * time.Millisecond)

	Convey("send a request without authentication", t, func() {
		Convey("should not be nil", func() {
			So(svr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(requestFlag, ShouldBeTrue)
			So(commonResp.Header.Code, ShouldEqual, errordef.Error_Msg_OperationFailed)
		})
	})
}
