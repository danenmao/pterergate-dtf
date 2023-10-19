package serverhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	Action    string `json:"Action"`
	RequestId string `json:"RequestId"`
	Version   string `json:"Version"`
	APIModule string `json:"ApiModule"`
	Sign      string `json:"Sign"`
	Timestamp string `json:"Timestamp"`
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

const ResponseBodyField string = "ResponseBody"

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
	returnErrorResponse(c, requestId, "InternalError", "internal error occurred")
}
