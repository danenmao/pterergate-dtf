package serverhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
)

// claim
type CommonMessage struct {
	UserName string `json:"username"`
	BodyHash string `json:"bodyhash"`
}

// Request
//
//	{
//	  "RequestHeader":{},
//	  "RequestBody":{}
//	}
type CommonRequest struct {
	Header RequestHeader `json:"RequestHeader"`
	Body   string        `json:"RequestBody"`
}

type RequestHeader struct {
	RequestId string `json:"RequestId"`
	Version   string `json:"Version"`
	BodyHash  string `json:"Sign"`
	Timestamp string `json:"Timestamp"`
	Module    string `json:"Module"`
	Action    string `json:"Action"`
}

// Response
//
//	{
//	 "ResponseHeader":{},
//	 "ResponseBody":{}
//	}
type CommonResponse struct {
	Header ResponseHeader `json:"ResponseHeader"`
	Body   string         `json:"ResponseBody"`
}

type ResponseHeader struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

type IRequest interface{}
type IResponse interface{}

func ReturnErrorResponse(requestId string, code string, message string) IResponse {
	return CommonResponse{
		Header: ResponseHeader{
			RequestId: requestId,
			Code:      code,
			Message:   message,
		},
	}
}

func returnErrorResponse(c *gin.Context, requestId string, code string, msg string) {
	response := ReturnErrorResponse(requestId, code, msg)
	c.JSON(http.StatusOK, response)
}

func returnInternalErrorResponse(c *gin.Context, requestId string) {
	returnErrorResponse(c, requestId, errordef.Error_Msg_InternalError, "internal error occurred")
}
