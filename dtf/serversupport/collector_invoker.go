package serversupport

import (
	"encoding/json"
	"fmt"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
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

type CollectorInvoker struct {
	InvokerBase
}

func NewCollectorInvoker(serverHost string, serverPort uint16, user string) *CollectorInvoker {
	return &CollectorInvoker{
		InvokerBase: *NewInvokerBase(serverHost, serverPort, user),
	}
}

// return an invoker function
// for executor to invoke collector
func (c *CollectorInvoker) GetInvoker() taskmodel.CollectorInvoker {
	return func(results []taskmodel.SubtaskResult) error {
		return c.invoker(results)
	}
}

func (c *CollectorInvoker) invoker(results []taskmodel.SubtaskResult) error {
	body := CollectorRequestBody{
		Results: results,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return errordef.ErrOperationFailed
	}

	return c.client.Post(c.url, c.UserName, string(data))
}
