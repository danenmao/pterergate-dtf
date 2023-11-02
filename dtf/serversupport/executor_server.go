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
	*ServerBase
	handler taskmodel.ExecutorRequestHandler
}

func NewExecutorServer(handler taskmodel.ExecutorRequestHandler) *ExecutorServer {
	return &ExecutorServer{
		ServerBase: NewServerBase(),
		handler:    handler,
	}
}

// for executor
// to register a handler to process executor requests
func (s *ExecutorServer) GetRegister() taskmodel.RegisterExecutorRequestHandler {
	return func(handler taskmodel.ExecutorRequestHandler) error {
		// save the handler
		s.handler = handler
		return nil
	}
}

// start the executor server to receive requests
func (s *ExecutorServer) Serve(serverPort uint16) error {
	return s.serve(
		serverPort, ExecutorServerURI,
		func(
			requestHeader serverhelper.RequestHeader,
			requestBody string,
		) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		},
	)
}

func (s *ExecutorServer) Shutdown() error {
	return s.server.Shutdown()
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

	err = s.handler(body.Subtasks)
	return "", err
}
