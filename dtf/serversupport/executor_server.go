package serversupport

import (
	"encoding/json"
	"errors"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/serverhelper"
)

const ExecutorServerURI = "/executor"

type ExecutorRequestBody struct {
	Subtasks []taskmodel.SubtaskBody `json:"Subtasks"`
}

type ExecutorServer struct {
	Handler taskmodel.ExecutorRequestHandler
	server  *serverhelper.SimpleServer
}

// for executor
// to register a handler to process executor requests
func (s *ExecutorServer) GetRegister() taskmodel.RegisterExecutorRequestHandler {
	return func(handler taskmodel.ExecutorRequestHandler) error {
		// save the handler
		s.Handler = handler
		return nil
	}
}

// start the executor server to receive requests
func (s *ExecutorServer) StartServer(uri string, serverPort uint16) error {
	s.server = serverhelper.NewSimpleServer(
		serverPort,
		map[string]serverhelper.RequestHandler{uri: func(
			requestHeader serverhelper.RequestHeader,
			requestBody string,
		) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		}},
	)

	s.server.StartServer()
	return nil
}

func (s *ExecutorServer) handleRequest(
	requestHeader serverhelper.RequestHeader,
	requestBody string,
) (response string, err error) {
	body := ExecutorRequestBody{}
	err = json.Unmarshal([]byte(requestBody), &body)
	if err != nil {
		return "", errors.New("failed to parse request body")
	}

	err = s.Handler(body.Subtasks)
	return "", err
}
