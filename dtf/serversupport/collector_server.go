package serversupport

import (
	"encoding/json"
	"errors"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/serverhelper"
)

const CollectorServerURI = "/collector"

// collector request body structure
type CollectorRequestBody struct {
	Results []taskmodel.SubtaskResult `json:"Results"`
}

// the Collector Server
// receive the requests of subtask results
type CollectorServer struct {
	*ServerBase
	handler taskmodel.CollectorRequestHandler
}

func NewCollectorServer(handler taskmodel.CollectorRequestHandler) *CollectorServer {
	return &CollectorServer{
		ServerBase: NewServerBase(),
		handler:    handler,
	}
}

// return a register function
// for collector
// to register a request handler to hander collector requests
func (s *CollectorServer) GetRegister() taskmodel.RegisterCollectorRequestHandler {
	return func(handler taskmodel.CollectorRequestHandler) error {
		// save the handler
		s.handler = handler
		return nil
	}
}

// start the collector server to receive requests
func (s *CollectorServer) Serve(serverPort uint16) error {
	return s.serve(
		serverPort, CollectorServerURI,
		func(
			requestHeader serverhelper.RequestHeader,
			requestBody string,
		) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		})
}

func (s *CollectorServer) handleRequest(
	requestHeader serverhelper.RequestHeader,
	requestBody string,
) (response string, err error) {
	body := CollectorRequestBody{}
	err = json.Unmarshal([]byte(requestBody), &body)
	if err != nil {
		return "", errors.New("failed to parse request body")
	}

	err = s.handler(body.Results)
	return "", err
}
