package serversupport

import (
	"encoding/json"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

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

	_, err = c.client.Post(c.url, c.UserName, string(data))
	return err
}
