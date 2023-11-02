package serversupport

import (
	"fmt"

	"github.com/danenmao/pterergate-dtf/internal/serverhelper"
)

type InvokerBase struct {
	ServerHost string
	ServerPort uint16
	URI        string
	UserName   string
	url        string
	client     *serverhelper.SimpleInvoker
}

func NewInvokerBase(serverHost string, serverPort uint16, user string) *InvokerBase {
	i := &InvokerBase{
		ServerHost: serverHost,
		ServerPort: serverPort,
		URI:        CollectorServerURI,
		UserName:   user,
		client:     serverhelper.NewSimpleInvoker(),
	}

	i.url = fmt.Sprintf("http://%s:%d%s", i.ServerHost, i.ServerPort, i.URI)
	return i
}
