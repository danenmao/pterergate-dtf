package serverhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/msgsigner"
)

const BODY_HASH = "BodyHash"
const USER_NAME = "UserName"

type RequestHandler func(header RequestHeader, requestBody string) (responseBody string, err error)
type SimpleServer struct {
	ServerPort uint16
	handlerMap map[string]RequestHandler
	server     *http.Server
	signer     *msgsigner.MsgSigner
}

func NewSimpleServer(serverPort uint16, handlerMap map[string]RequestHandler) *SimpleServer {
	server := &SimpleServer{
		handlerMap: make(map[string]RequestHandler),
		ServerPort: serverPort,
		server:     nil,
		signer:     msgsigner.NewMsgSigner(),
	}

	if handlerMap != nil {
		server.handlerMap = handlerMap
	}

	return server
}

func (s *SimpleServer) Serve() error {
	ginMode := gin.DebugMode
	if config.WorkEnv == config.ENV_ONLINE {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	router := gin.Default()
	for uri, handler := range s.handlerMap {
		router.POST(
			uri,
			s.requestTracing(),
			s.authMiddleware(),
			s.handleRequest(handler),
		)
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.ServerPort),
		Handler: router,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		glog.Error("failed to run the gin server: ", err.Error())
	}

	glog.Info("exited")
	return nil
}

func (s *SimpleServer) Shutdown() error {
	return s.ShutdownWithDuration(5 * time.Second)
}

func (s *SimpleServer) ShutdownWithDuration(duration time.Duration) error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		glog.Error("The server shutdown error: ", err)
		return err
	}

	glog.Info("The server shutdown successfully.")
	return nil
}

func (s *SimpleServer) requestTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (s *SimpleServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Info("verify the Authentication header")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed,
				"No Authorization")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed,
				"invalid Authorization format")
			c.Abort()
			return
		}

		msg, err := s.signer.Verify(parts[1])
		if err != nil {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed,
				"invalid Authorization token")
			c.Abort()
			return
		}

		commonMsg := CommonMessage{}
		err = json.Unmarshal([]byte(msg), &commonMsg)
		if err != nil {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed,
				"invalid token data")
			c.Abort()
			return
		}

		c.Set(USER_NAME, commonMsg.UserName)
		c.Set(BODY_HASH, commonMsg.BodyHash)
		c.Next()
	}
}

func (s *SimpleServer) handleRequest(handler RequestHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			glog.Warning("failed to get body: ", err.Error())
			returnErrorResponse(c, "", errordef.Error_Msg_ParsingParam,
				"NO request id found")
			return
		}

		var request = CommonRequest{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			glog.Warning("failed to parse common parameter: ", err.Error())
			glog.Warning(string(body))
			returnErrorResponse(c, "", errordef.Error_Msg_ParsingParam,
				"failed to parse parameter")
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		actualBodyHash := CalcMsgHash(request.Body)
		expectedHash, existed := c.Get(BODY_HASH)
		if !existed {
			returnInternalErrorResponse(c, request.Header.RequestId)
			return
		}

		if actualBodyHash != expectedHash {
			returnErrorResponse(c, request.Header.RequestId,
				errordef.Error_Msg_AuthorizationFailed,
				"invalid body hash")
			return
		}

		start := time.Now()
		response, err := s.invokeHandler(handler, request)
		timeCost := time.Since(start)
		glog.Info("handler stat, requestId:", request.Header.RequestId, ", timeCost: ", timeCost)

		if err == nil {
			c.JSON(http.StatusOK, response)
		} else {
			returnInternalErrorResponse(c, request.Header.RequestId)
		}
	}
}

func (s *SimpleServer) invokeHandler(
	handler RequestHandler,
	request CommonRequest,
) (response IResponse, err error) {
	// invoke the outer handler
	rspBody, err := handler(request.Header, request.Body)
	if err != nil {
		return ReturnErrorResponse(
			request.Header.RequestId,
			errordef.Error_Msg_OperationFailed,
			err.Error()), nil
	}

	// return response
	response = CommonResponse{
		Header: ResponseHeader{
			RequestId: request.Header.RequestId,
			Code:      errordef.Error_Msg_Success,
			Message:   errordef.Error_Msg_Success,
		},
		Body: rspBody,
	}

	return
}
