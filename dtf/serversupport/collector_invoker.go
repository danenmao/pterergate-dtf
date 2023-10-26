package serversupport

import (
	"encoding/json"
	"fmt"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

type CollectorInvoker struct {
	ServerHost string
	ServerPort uint16
	URI        string
	UserName   string
	url        string
	client     *SimpleInvoker
}

func NewCollectorInvoker(serverHost string, serverPort uint16, user string) *CollectorInvoker {
	c := &CollectorInvoker{
		ServerHost: serverHost,
		ServerPort: serverPort,
		URI:        CollectorServerURI,
		UserName:   user,
		client:     NewSimpleInvoker(),
	}

	c.url = fmt.Sprintf("http://%s:%d%s", c.ServerHost, c.ServerPort, c.URI)
	return nil
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
