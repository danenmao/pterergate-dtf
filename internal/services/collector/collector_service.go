package collector

import "github.com/danenmao/pterergate-dtf/dtf/taskmodel"

// handle collector requests
func CollectorRequestHandler(results []taskmodel.SubtaskResult) error {
	InsertSubtaskResults(results)
	return nil
}
