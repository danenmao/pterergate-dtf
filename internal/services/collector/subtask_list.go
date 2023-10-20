package collector

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

type SubtaskElem struct {
	Result     *taskmodel.SubtaskResult
	InsertTime time.Time
}

const (
	MaxSubtaskElemCount = 10000
	MaxCountPerTime     = 10
)

var s_SubtaskElemList = []*SubtaskElem{}
var s_SubtaskLock sync.RWMutex

func InsertSubtask(result *taskmodel.SubtaskResult) {
	s_SubtaskLock.Lock()
	defer s_SubtaskLock.Unlock()

	// check if exceed the max length
	if len(s_SubtaskElemList) >= MaxSubtaskElemCount {
		glog.Info("the length of subtask list exceeded the max count")
		return
	}

	// append to the list
	subtask := SubtaskElem{Result: result, InsertTime: time.Now()}
	s_SubtaskElemList = append(s_SubtaskElemList, &subtask)
}

func InsertSubtaskResults(results []taskmodel.SubtaskResult) {
	s_SubtaskLock.Lock()
	defer s_SubtaskLock.Unlock()

	cap := MaxSubtaskElemCount - len(s_SubtaskElemList)
	count := len(results)
	if count > cap {
		count = cap
	}

	if count <= 0 {
		glog.Info("the length of subtask list exceeded the max count")
		return
	}

	for _, result := range results[0:count] {
		subtask := SubtaskElem{Result: &result, InsertTime: time.Now()}
		s_SubtaskElemList = append(s_SubtaskElemList, &subtask)
	}
}

func InsertSubtaskList(newElems *[]*SubtaskElem) {
	s_SubtaskLock.Lock()
	defer s_SubtaskLock.Unlock()

	if len(s_SubtaskElemList)+len(*newElems) >= MaxSubtaskElemCount {
		glog.Warning("too many elements in the list")
		return
	}

	s_SubtaskElemList = append(s_SubtaskElemList, *newElems...)
}

func PopSubtaskList(retList *[]*SubtaskElem) {

	s_SubtaskLock.Lock()
	defer s_SubtaskLock.Unlock()

	totalCount := uint(len(s_SubtaskElemList))
	if totalCount <= 0 {
		return
	}
	glog.Info("subtask elem count in list: ", totalCount)

	// get the element count
	count := 1
	if totalCount > uint(gs_RoutineLimit.UpperLimit) {
		count = int(totalCount/uint(gs_RoutineLimit.UpperLimit)) + 1
	}

	// MaxCountPerTime
	if count > MaxCountPerTime {
		count = MaxCountPerTime
	}

	getList := s_SubtaskElemList[0:count]
	*retList = append(*retList, getList...)

	// remove the elements
	s_SubtaskElemList = s_SubtaskElemList[count:]
}
