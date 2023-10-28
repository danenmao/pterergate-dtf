package serverhelper

import (
	"bytes"
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
	Handler    RequestHandler
	URI        string
	ServerPort uint16
	signer     *msgsigner.MsgSigner
}

func NewSimpleServer(uri string, serverPort uint16, handler RequestHandler) *SimpleServer {
	return &SimpleServer{
		URI:        uri,
		ServerPort: serverPort,
		Handler:    handler,
		signer:     msgsigner.NewMsgSigner(),
	}
}

func (s *SimpleServer) StartServer() error {
	ginMode := gin.DebugMode
	if config.WorkEnv == config.ENV_ONLINE {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	router := gin.Default()
	router.POST(
		s.URI,
		s.requestTracing(),
		s.authMiddleware(),
		s.handleCommonRequest(),
	)

	err := router.Run(fmt.Sprintf(":%d", s.ServerPort))
	if err != nil {
		glog.Error("failed to run gin: ", err.Error())
	}

	glog.Info("exited")
	return nil
}

func (s *SimpleServer) requestTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (s *SimpleServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Info("Authentication")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed, "No Authorization")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed, "invalid Authorization format")
			c.Abort()
			return
		}

		msg, err := s.signer.Verify(parts[1])
		if err != nil {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed, "invalid Authorization token")
			c.Abort()
			return
		}

		commonMsg := CommonMessage{}
		err = json.Unmarshal([]byte(msg), &commonMsg)
		if err != nil {
			returnErrorResponse(c, "", errordef.Error_Msg_AuthorizationFailed, "invalid token data")
			c.Abort()
			return
		}

		c.Set(USER_NAME, commonMsg.UserName)
		c.Set(BODY_HASH, commonMsg.BodyHash)
		c.Next()
	}
}

func (s *SimpleServer) handleCommonRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			glog.Warning("failed to get body: ", err.Error())
			returnErrorResponse(c, "", errordef.Error_Msg_ParsingParam, "NO request id found")
			return
		}

		var request = CommonRequest{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			glog.Warning("failed to parse common parameter: ", err.Error())
			glog.Warning(string(body))
			returnErrorResponse(c, "", errordef.Error_Msg_ParsingParam, "failed to parse parameter")
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
			returnErrorResponse(c, request.Header.RequestId, errordef.Error_Msg_AuthorizationFailed,
				"invalid body hash")
			return
		}

		start := time.Now()
		response, err := s.handle(request)
		timeCost := time.Since(start)
		glog.Info("action stat, requestId:", request.Header.RequestId, ", timeCost: ", timeCost)

		if err == nil {
			c.JSON(http.StatusOK, response)
		} else {
			returnInternalErrorResponse(c, request.Header.RequestId)
		}
	}
}

func (s *SimpleServer) handle(request CommonRequest) (response IResponse, err error) {
	// invoke the outer handler
	rspBody, err := s.Handler(request.Header, request.Body)
	if err != nil {
		return ReturnErrorResponse(request.Header.RequestId,
			errordef.Error_Msg_OperationFailed, err.Error()), nil
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
