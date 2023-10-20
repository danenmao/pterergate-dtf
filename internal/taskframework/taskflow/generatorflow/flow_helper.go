package generatorflow

import (
	"errors"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/subtaskqueue"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

const (
	SubtaskGenerationInterval = 10
	SubtaskGenerationMaxTime  = 3600
)

// generator flow helpr
type GeneratorFlowHelper struct {
	SubtaskQueues subtaskqueue.SubtaskQueueMgr                // subtask queue manager
	GeneratorMap  map[taskmodel.TaskIdType]TaskGenerationImpl // task generator object map
	Mutex         sync.Mutex                                  // lock
}

type TaskGenerationImpl struct {
	Impl     taskmodel.ITaskGenerator
	TaskType uint32
}

func NewGeneratorFlowHelper() GeneratorFlowHelper {
	return GeneratorFlowHelper{
		SubtaskQueues: subtaskqueue.SubtaskQueueMgr{
			SubtaskQueueMap: make(map[taskmodel.TaskIdType]*subtaskqueue.SubtaskQueue),
		},
		GeneratorMap: make(map[taskmodel.TaskIdType]TaskGenerationImpl),
		Mutex:        sync.Mutex{},
	}
}

// global generator helper object
var gs_GeneratorHelper = NewGeneratorFlowHelper()

func GetGeneratorFlowHelper() *GeneratorFlowHelper {
	return &gs_GeneratorHelper
}

func (generator *GeneratorFlowHelper) Begin(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	taskData *taskmodel.TaskParam,
	taskGenerator taskmodel.ITaskGenerator,
) error {

	generator.Mutex.Lock()
	defer generator.Mutex.Unlock()

	// create a subtask queue of the task
	generator.SubtaskQueues.AddTask(taskId)

	// save into the task generator map
	generator.GeneratorMap[taskId] = TaskGenerationImpl{
		Impl:     taskGenerator,
		TaskType: taskType,
	}

	return nil
}

func (generator *GeneratorFlowHelper) End(
	taskId taskmodel.TaskIdType,
) error {

	generator.Mutex.Lock()
	defer generator.Mutex.Unlock()

	// delete the subtask queue of the task
	err := generator.SubtaskQueues.RemoveTask(taskId)
	if err != nil {
		glog.Warning("failed to remove task: ", taskId, ", ", err.Error())
	}

	delete(generator.GeneratorMap, taskId)

	return nil
}

func (generator *GeneratorFlowHelper) GenerationLoop(
	taskId taskmodel.TaskIdType,
) error {

	// search the generator of the task
	generator.Mutex.Lock()
	impl, ok := generator.GeneratorMap[taskId]
	generator.Mutex.Unlock()
	if !ok {
		glog.Warning("task id not found in subtask fn map: ", taskId)
		return errors.New("task id not found")
	}

	// 执行子任务生成操作
	return generator.pickupSubtaskLoop(taskId, &impl)
}

func (generator *GeneratorFlowHelper) pickupSubtaskLoop(
	taskId taskmodel.TaskIdType,
	impl *TaskGenerationImpl,
) error {

	// record the start time
	startTime := time.Now().Unix()
	renewTime := startTime

	// create a routine to refresh the status
	exitChan := make(chan bool, 1)
	go asyncRefreshGenerationStatus(taskId, exitChan)

	// generation loop
	for {

		// try to create a subtask from the plugin generator
		finished := false
		subtaskData := taskmodel.SubtaskBody{}
		err := CreateSubtask(taskId, impl.TaskType, impl.Impl, &subtaskData, &finished)
		if err != nil && err != errordef.ErrNotFound {
			glog.Warning("failed to create a subtask: ", taskId, ",", err.Error())
			break
		}

		// check if get a subtask
		gotSubtask := (err == nil)

		// periodically save the task generation status
		taskStatus := ""
		taskStatus, err = impl.Impl.SaveStatus(taskId)
		if err != nil {
			glog.Warning("failed to save task status: ", taskId, ", ", err.Error())
		} else {
			err = SaveStatus(taskId, taskStatus)
			if err != nil {
				glog.Warning("failed to save task status: ", taskId, ", ", err.Error())
			}
		}

		// push the subtask into the subtask queue
		if gotSubtask {
			err = generator.SubtaskQueues.PushSubtask(taskId, &subtaskData)
			if err != nil {
				glog.Warning("failed to push subtask: ", taskId, ", ", err.Error())
				break
			}
		}

		// the task generation is over
		if finished {
			glog.Info("generation loop finished, break: ", taskId)
			break
		}

		// control the max generation time cost
		endTime := time.Now().Unix()
		if endTime-startTime >= SubtaskGenerationMaxTime {
			glog.Warning("task generation exceeds max generation time: ", taskId)
			break
		}

		// control the generation interval and speed
		time.Sleep(time.Millisecond * SubtaskGenerationInterval)

		// renew the generation ownership
		if endTime-renewTime >= 5 {
			tasktool.RenewTask(taskId)
			renewTime = endTime
			tasktool.UpdateTaskGenerationNextCheckTime(taskId)
		}
	}

	// notify to exit
	exitChan <- true
	close(exitChan)

	return nil
}

// refresh the generation status
func asyncRefreshGenerationStatus(
	taskId taskmodel.TaskIdType,
	exitChan chan bool,
) {

	for {

		// check if generation completed
		select {
		case <-exitChan:
			glog.Info("exit task generator refresh routine: ", taskId)
			return
		default:
			glog.Info("to refresh task generator status: ", taskId)
		}

		tasktool.RenewTask(taskId)
		tasktool.UpdateTaskGenerationNextCheckTime(taskId)

		time.Sleep(time.Second * 30)
	}
}
