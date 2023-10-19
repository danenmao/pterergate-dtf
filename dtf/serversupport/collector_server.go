package serversupport

import (
	"github.com/danenmao/pterergate-dtf/dtf/serversupport/serverhelper"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

// collector request body structure
type CollectorRequestBody struct {
	Results []taskmodel.SubtaskResult `json:"Results"`
}

// the Collector Server
type CollectorServer struct {
	Handler    taskmodel.CollectorRequestHandler
	URI        string
	ServerPort uint16
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
func (s *CollectorServer) StartServer() error {
	server := serverhelper.Server{
		URI:        s.URI,
		ServerPort: s.ServerPort,
		Handler: func(requestHeader serverhelper.RequestHeader, requestBody string) (response string, err error) {
			return s.handleRequest(requestHeader, requestBody)
		},
	}

	server.StartServer()
	return nil
}

func (s *CollectorServer) handleRequest(requestHeader serverhelper.RequestHeader, requestBody string) (response string, err error) {
	return "", nil
}
