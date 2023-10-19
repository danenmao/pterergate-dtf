package serversupport

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

// for scheduler to invoke the executor
// executor invoker, send subtasks to the executor
func ExecutorInvoker(subtaskData []taskmodel.SubtaskBody) error {
	// Host:Port
	// POST /executor
	return nil
}

// for executor
// to register a handler to process executor requests
func RegisterExecutorRequestHandler(handler taskmodel.ExecutorRequestHandler) error {
	// save the handler
	return nil
}
