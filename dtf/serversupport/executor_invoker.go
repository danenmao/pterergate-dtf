package serversupport

import (
	"encoding/json"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

type ExecutorInvoker struct {
	InvokerBase
}

func NewExecutorInvoker(serverHost string, serverPort uint16, user string) *ExecutorInvoker {
	return &ExecutorInvoker{
		InvokerBase: *NewInvokerBase(serverHost, serverPort, user),
	}
}

// return an invoker function
// for scheduler to invoke the executor
// executor invoker, send subtasks to the executor
func (e *ExecutorInvoker) GetInvoker() taskmodel.ExecutorInvoker {
	return func(subtaskBody []taskmodel.SubtaskBody) error {
		return e.invoker(subtaskBody)
	}
}

func (e *ExecutorInvoker) invoker(subtasks []taskmodel.SubtaskBody) error {
	body := ExecutorRequestBody{
		Subtasks: subtasks,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return errordef.ErrOperationFailed
	}

	return e.client.Post(e.url, e.UserName, string(data))
}
