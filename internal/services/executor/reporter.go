package executor

import (
	"sync"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
)

type SubtaskResultReporter struct {
	Results []taskmodel.SubtaskResult
	Lock    sync.Mutex
}

const (
	MaxSubtaskElemCount = 10000
	MaxCountPerTime     = 10
)

var gs_ResultReporter SubtaskResultReporter

func GetReporter() *SubtaskResultReporter {
	return &gs_ResultReporter
}

func ReportRoutine() {

	// get results
	results := []taskmodel.SubtaskResult{}
	GetReporter().PopSubtaskResult(&results)

	// send to collector
	GetReporter().ReportToCollector(results)
}

func (reporter *SubtaskResultReporter) AddSubtaskResult(result *taskmodel.SubtaskResult) error {
	reporter.Lock.Lock()
	defer reporter.Lock.Unlock()

	// too many results, refuse
	if len(reporter.Results) >= MaxSubtaskElemCount {
		glog.Warning("the length of result list exceeded the max count")
		return nil
	}

	// add to reporter list
	reporter.Results = append(reporter.Results, *result)
	return nil
}

func (reporter *SubtaskResultReporter) PopSubtaskResult(retList *[]taskmodel.SubtaskResult) error {
	reporter.Lock.Lock()
	defer reporter.Lock.Unlock()

	count := len(reporter.Results)
	if count > MaxCountPerTime {
		count = MaxCountPerTime
	}

	getList := reporter.Results[0:count]
	*retList = append(*retList, getList...)
	reporter.Results = reporter.Results[count:]
	return nil
}

func (reporter *SubtaskResultReporter) ReportToCollector(results []taskmodel.SubtaskResult) {
	CollectorInvoker(results)
}
