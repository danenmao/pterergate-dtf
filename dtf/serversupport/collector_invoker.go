package serversupport

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

type CollectorInvoker struct {
	ServerHost string
	ServerPort uint16
	URI        string
}

// return an invoker function
// for executor to invoke collector
func (c *CollectorInvoker) GetInvoker() taskmodel.CollectorInvoker {
	return func(results []taskmodel.SubtaskResult) error {
		// Host:Port
		// POST /collector
		return nil
	}
}
