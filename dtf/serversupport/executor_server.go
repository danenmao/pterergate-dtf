package serversupport

import (
	"github.com/danenmao/pterergate-dtf/dtf/serversupport/serverhelper"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

const ExecutorServerURI = "/executor"

type ExecutorRequestBody struct {
	Subtasks []taskmodel.SubtaskBody `json:"Subtasks"`
}

type ExecutorServer struct {
	Handler    taskmodel.ExecutorRequestHandler
	URI        string
	ServerPort uint16
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
func (s *ExecutorServer) StartServer() error {
	server := serverhelper.SimpleServer{
		URI:        s.URI,
		ServerPort: s.ServerPort,
		Handler: func(
			requestHeader serverhelper.RequestHeader,
			requestBody string,
		) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		},
	}

	server.StartServer()
	return nil
}

func (s *ExecutorServer) handleRequest(
	requestHeader serverhelper.RequestHeader,
	requestBody string,
) (response string, err error) {
	return "", nil
}
