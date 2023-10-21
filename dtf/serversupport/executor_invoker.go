package serversupport

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

type ExecutorInvoker struct {
	ServerHost string
	ServerPort uint16
	URI        string
}

// return an invoker function
// for scheduler to invoke the executor
// executor invoker, send subtasks to the executor
func (e *ExecutorInvoker) GetInvoker() taskmodel.ExecutorInvoker {
	return func(subtaskBody []taskmodel.SubtaskBody) error {
		// Host:Port
		// POST /collector
		return nil
	}
}
