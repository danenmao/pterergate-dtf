package serversupport

import (
	"errors"

	"github.com/danenmao/pterergate-dtf/internal/exitctrl"
	"github.com/danenmao/pterergate-dtf/internal/serverhelper"
)

type ServerBase struct {
	server *serverhelper.SimpleServer
}

func NewServerBase() *ServerBase {
	return &ServerBase{}
}

func (s *ServerBase) serve(serverPort uint16, uri string, handler serverhelper.RequestHandler) error {
	s.server = serverhelper.NewSimpleServer(
		serverPort,
		map[string]serverhelper.RequestHandler{uri: handler},
	)

	if s.server == nil {
		return errors.New("failed to new a simple server")
	}

	exitctrl.AddExitRoutine(func() {
		s.server.Shutdown()
	})

	return s.server.Serve()
}

func (s *ServerBase) Shutdown() error {
	return s.server.Shutdown()
}
