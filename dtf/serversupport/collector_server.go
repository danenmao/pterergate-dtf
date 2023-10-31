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
	Handler taskmodel.CollectorRequestHandler
	server  *serverhelper.SimpleServer
}

// return a register function
// for collector
// to register a request handler to hander collector requests
func (s *CollectorServer) GetRegister() taskmodel.RegisterCollectorRequestHandler {
	return func(handler taskmodel.CollectorRequestHandler) error {
		// save the handler
		s.Handler = handler
		return nil
	}
}

// start the collector server to receive requests
func (s *CollectorServer) StartServer(uri string, serverPort uint16) error {
	s.server = serverhelper.NewSimpleServer(
		uri, serverPort,
		func(
			requestHeader serverhelper.RequestHeader,
			requestBody string,
		) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		},
	)

	s.server.StartServer()
	return nil
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

	err = s.Handler(body.Results)
	return "", err
}
