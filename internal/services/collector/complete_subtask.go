package collector

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/subtasktool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/collectorlogic"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

const (
	// complete subtask routine limit
	SubtaskRoutineCountDefaultLimit = 300
	SubtaskRoutineLimitEnvName      = "SUBTASK_ROUTINE_LIMIT"
)

var (
	// control the complete subtask routine limit
	gs_RoutineLimit = routine.CountLimiter{
		UpperLimit: SubtaskRoutineCountDefaultLimit,
	}
)

// to complete subtask
func CompleteSubtaskRoutine() {

	startTime := time.Now()
	for {
		checkCompletedSubtask()
		if time.Since(startTime) > time.Millisecond*900 {
			break
		}

		time.Sleep(time.Millisecond)
	}

}

func checkCompletedSubtask() {

	// check if create a new complete subtask routine
	if !gs_RoutineLimit.IncrIfNotFull() {
		return
	}

	toDecr := true
	defer func() {
		if toDecr {
			gs_RoutineLimit.Decr()
		}
	}()

	// get subtasks list
	var list = []*SubtaskElem{}
	PopSubtaskList(&list)
	if len(list) == 0 {
		return
	}

	toDecr = false
	go completeSubtask(list)
}

func completeSubtask(list []*SubtaskElem) {

	// decr the ref
	defer gs_RoutineLimit.Decr()

	// stat the time cost
	startTime := time.Now()
	defer func() {
		cost := time.Since(startTime)
		glog.Info("processing subtasks costs: ", cost)
	}()

	// to complete the subtasks
	err := doCompleteSubtask(list)
	if err != nil {
		// if failed, push them back to the list
		glog.Warning("re-insert elements into the list: ", err, len(list))
		InsertSubtaskList(&list)
	}
}

func doCompleteSubtask(elems []*SubtaskElem) error {

	var zlist = []*redis.Z{}
	var idList = []interface{}{}
	pipeline := redistool.DefaultRedis().Pipeline()
	endTime := time.Now().Unix()

	glog.Info("subtask count to process: ", len(elems))
	for _, elem := range elems {

		if elem == nil {
			panic("invalid elem pointer")
		}

		if elem.Result == nil {
			glog.Error("invalid body pointer")
			continue
		}

		result := elem.Result
		glog.Info("begin to process subtask: ", result.SubtaskId, result.TaskId)

		subtaskCompleted := false
		err := processSubtaskResult(result, pipeline, &subtaskCompleted)
		if err != errordef.ErrNotFound && err != nil {
			continue
		}

		// not completed, skip
		if !subtaskCompleted {
			continue
		}

		z := redis.Z{
			Score:  float64(endTime),
			Member: result.SubtaskId,
		}

		zlist = append(zlist, &z)
		idList = append(idList, result.SubtaskId)
		glog.Infof("processed completed subtask: ", result.SubtaskId)
	} // for

	// remove from running subtask list
	if len(idList) > 0 {
		pipeline.ZRem(context.Background(), config.RunningSubtaskZset, idList...)
	}

	// insert to completed subtask list
	if len(zlist) > 0 {
		pipeline.ZAdd(context.Background(), config.CompletedSubtaskList, zlist...)
	}

	// exec pipeline
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to add subtask to redis: ", err, len(elems))
		return err
	}

	glog.Info("succeed to process completed subtask: ", idList)
	return nil
}

func processSubtaskResult(
	result *taskmodel.SubtaskResult,
	pipeline redis.Pipeliner,
	subtaskCompleted *bool,
) error {

	*subtaskCompleted = false

	running := subtasktool.IsSubtaskRunning(result.SubtaskId)
	if !running {
		glog.Info("subtask is not running: ", result.SubtaskId)
		return nil
	}

	running = tasktool.IsTaskRunning(result.TaskId)
	if !running {
		glog.Info("task is not running: ", result.TaskId)
		return nil
	}

	collectorlogic.OnSubtaskResult(result, subtaskCompleted)
	if *subtaskCompleted {
		SetSubtaskResult(result.SubtaskId, result, &pipeline)
	}

	return nil
}

func SetSubtaskResult(
	subtaskId taskmodel.SubtaskIdType,
	result *taskmodel.SubtaskResult,
	ppipeline *redis.Pipeliner,
) error {

	data, err := json.Marshal(result)
	if err != nil {
		glog.Warning("failed to serialize subtask: ", result.SubtaskId, result.TaskId)
		data = []byte("")
	}

	err = subtasktool.SetSubtaskResult(uint64(subtaskId), result.Result, string(data), ppipeline)
	if err != nil {
		glog.Warning("failed to set subtask result: ", subtaskId, err)
		return err
	}

	return nil
}
