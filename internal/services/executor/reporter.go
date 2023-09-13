package executor

import (
	"sync"

	"pterergate-dtf/dtf/taskmodel"
)

type SubtaskResultReporter struct {
	Results []taskmodel.SubtaskResult
	Lock    sync.Mutex
}

var gs_ResultReporter SubtaskResultReporter

func GetReporter() *SubtaskResultReporter {
	return &gs_ResultReporter
}

func (reporter *SubtaskResultReporter) AddSubtaskResult(result *taskmodel.SubtaskResult) error {
	return nil
}

func ReportRoutine() {

}

func (reporter *SubtaskResultReporter) ReportCollector(results []taskmodel.SubtaskResult) {
	CollectorInvoker(results)
}
