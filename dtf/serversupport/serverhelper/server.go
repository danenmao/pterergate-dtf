package serverhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/config"
)

type RequestHandler func(header RequestHeader, requestBody string) (responseBody string, err error)
type Server struct {
	Handler    RequestHandler
	URI        string
	ServerPort uint16
}

func (s *Server) StartServer() error {
	ginMode := gin.DebugMode
	if config.WorkEnv == config.ENV_ONLINE {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	router := gin.Default()
	router.POST(
		s.URI,
		s.RequestTracing(),
		s.AuthMiddleware(),
		s.HandleCommonRequest(),
	)

	err := router.Run(fmt.Sprintf(":%d", s.ServerPort))
	if err != nil {
		glog.Error("failed to run gin: ", err.Error())
	}

	glog.Info("exited")
	return nil
}

func (s *Server) RequestTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Info("Authentication")
	}
}

func (s *Server) HandleCommonRequest() gin.HandlerFunc {
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

		commonParam := request.Header
		actionName := commonParam.Action
		requestId := commonParam.RequestId

		start := time.Now()
		response, err := s.processByAction(request)
		timeCost := time.Since(start)
		glog.Info("action stat, action: ", actionName, ", requestId: ", requestId, ", timeCost: ", timeCost)

		if err == nil {
			c.JSON(http.StatusOK, response)
		} else {
			returnInternalErrorResponse(c, requestId)
		}
	}
}

func (s *Server) processByAction(request CommonRequest) (response IResponse, err error) {

	s.Handler(request.Header, "")
	return ReturnErrorResponse(
		request.Header.RequestId,
		errordef.Error_Msg_InvalidParameter,
		"unknown action"), nil
}
