package executor

import (
	"sync"
	"time"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/dtf/taskplugin"
	"pterergate-dtf/internal/taskframework/taskloader"
)

// collector service invoker
var CollectorInvoker taskmodel.CollectorInvoker

type ExecutorService struct {
	ExecutorMap map[uint32]taskmodel.ITaskExecutor
	Lock        sync.Mutex
}

var gs_ExecutorService ExecutorService

func GetExecutorService() *ExecutorService {
	return &gs_ExecutorService
}

// handler executor service request
func ExecutorRequestHandler(subtasks []taskmodel.SubtaskData) error {

	// check if exceed the subtask count

	// execute each subtask in a go routine
	for _, subtask := range subtasks {
		go GetExecutorService().execSubtask(&subtask)
	}
	return nil
}

func (service *ExecutorService) Init() error {
	service.ExecutorMap = map[uint32]taskmodel.ITaskExecutor{}
	return nil
}

// to execute subtask
func (service *ExecutorService) execSubtask(subtask *taskmodel.SubtaskData) error {

	// get the task executor object
	var executor taskmodel.ITaskExecutor
	err := service.getTaskExecutor(subtask.TaskType, &executor)
	if err != nil {
		glog.Warning("failed to get the task executor: ", err)
		return err
	}

	// execute this subtask asynchronously
	resultChan := make(chan taskmodel.SubtaskResult, 1)
	go func() {
		result := taskmodel.SubtaskResult{
			TaskId:    subtask.TaskId,
			SubtaskId: subtask.SubtaskId,
		}

		err := executor.Execute(subtask, &result)
		if err != nil {
			glog.Warning("TaskExecutor returned err: ", err)
			result.Result = taskmodel.SubtaskResult_Failure
			result.Reason = err.Error()
		} else {
			result.Result = taskmodel.SubtaskResult_Success
			result.Reason = "success"
		}

		resultChan <- result
	}()

	result := taskmodel.SubtaskResult{
		TaskId:    subtask.TaskId,
		SubtaskId: subtask.SubtaskId,
	}
	select {
	case result = <-resultChan:
		glog.Info("subtask completed: ", subtask.SubtaskId, " of ", subtask.TaskId)

	case <-time.After(time.Second * time.Duration(subtask.Timeout)):
		glog.Info("subtask timeout: ", subtask.SubtaskId, " of ", subtask.TaskId)
		result.Result = taskmodel.SubtaskResult_Timeout
		result.Reason = "timeout"
	}

	subtask.TerminatedAt = time.Now()

	// add result to notify queue
	GetReporter().AddSubtaskResult(&result)
	return nil
}

func (service *ExecutorService) getTaskExecutor(taskType uint32, retExecutor *taskmodel.ITaskExecutor) error {

	service.Lock.Lock()
	executor, ok := service.ExecutorMap[taskType]
	service.Lock.Unlock()
	if ok {
		*retExecutor = executor
		return nil
	}

	err := GetTaskExecutor(taskType, retExecutor)
	if err != nil {
		return err
	}

	service.Lock.Lock()
	service.ExecutorMap[taskType] = *retExecutor
	service.Lock.Unlock()
	return nil
}

func GetTaskExecutor(taskType uint32, executor *taskmodel.ITaskExecutor) error {

	var plugin taskplugin.ITaskPlugin = nil
	err := taskloader.LookupTaskPlugin(taskType, &plugin)
	if err != nil {
		glog.Warning("failed to get task plugin: ", taskType)
		return err
	}

	var taskBody taskmodel.TaskBody
	err = plugin.GetTaskBody(&taskBody)
	if err != nil {
		glog.Warning("failed to get task context: ", err.Error())
		return err
	}

	*executor = taskBody.Executor
	glog.Info("succeeded to get task scheduler: ", taskType)
	return nil
}
