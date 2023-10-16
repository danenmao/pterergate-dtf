package collector

import "pterergate-dtf/dtf/taskmodel"

// handle collector requests
func CollectorRequestHandler(results []taskmodel.SubtaskResult) error {
	InsertSubtaskResults(results)
	return nil
}
