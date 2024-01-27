package generationlogic

import (
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/dtf/taskplugin"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskloader"
)

// generation flow
type GenerationFlow struct {
	TaskId    taskmodel.TaskIdType
	TaskType  uint32
	Generator taskmodel.ITaskGenerator
	TaskParam *taskmodel.TaskParam
}

// create a generation flow object
// every task has their own generation flow object
func NewGenerationFlow() *GenerationFlow {
	return &GenerationFlow{}
}

// init the generation
func (flow *GenerationFlow) InitGeneration(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	taskParam *taskmodel.TaskParam,
) error {
	flow.TaskId = taskId
	flow.TaskType = taskType
	flow.TaskParam = taskParam

	// get the generator instance of this task type
	err := GetTaskGenerator(taskType, &flow.Generator)
	if err != nil {
		glog.Warning("failed to get task generator: ", taskId, taskType, ",", err)
		return err
	}

	err = GetFlowHelper().Begin(taskId, taskType, taskParam, flow.Generator)
	if err != nil {
		glog.Warning("failed to invoke GeneratorFlowHelper.Begin: ", taskId, ", ", taskType, ", ", err)
		return err
	}

	// try to load fomrer generation status of this task
	lastStatus := ""
	err = LoadStatus(taskId, &lastStatus)
	if err != nil {
		glog.Warning("failed to load task status: ", taskId, ", ", err.Error())
		return err
	}

	// begin to generate
	err = flow.Generator.Begin(taskId, taskType, taskParam, lastStatus)
	if err != nil {
		glog.Warning("generator.Begin failed: ", taskId, ",", err)
		return err
	}

	glog.Info("succeeded to init task generation: ", taskId)
	return nil
}

// finish the generation process
func (flow *GenerationFlow) FinishGeneration() error {
	// invoke the generator
	err := flow.Generator.End(flow.TaskId)
	if err != nil {
		glog.Warning("failed to finish the generation: ", flow.TaskId)
		return err
	}

	err = GetFlowHelper().End(flow.TaskId)
	if err != nil {
		glog.Warning("failed to invoke GeneratorFlowHelper.End: ", flow.TaskId, ", ", err)
		return err
	}

	glog.Info("succeeded to finish task generation: ", flow.TaskId)
	return nil
}

// generation loop
func (flow *GenerationFlow) GenerationLoop() error {
	err := GetFlowHelper().GenerationLoop(flow.TaskId)
	if err != nil {
		glog.Warning("GeneratorFlowHelper.GenerationLoop failed: ", flow.TaskId, ", ", err)
		return err
	}

	glog.Info("GeneratorHelper.GenerationLoop succeeded: ", flow.TaskId)
	return nil
}

// get the generator object of the specified task type
func GetTaskGenerator(taskType uint32, generator *taskmodel.ITaskGenerator) error {
	var plugin taskplugin.ITaskPlugin = nil
	err := taskloader.LookupTaskPlugin(taskType, &plugin)
	if err != nil {
		glog.Warning("failed to get task plugin: ", taskType)
		return err
	}

	var body taskmodel.PluginBody
	err = plugin.GetPluginBody(&body)
	if err != nil {
		glog.Warning("failed to get task context: ", err.Error())
		return err
	}

	*generator = body.Generator
	glog.Info("succeeded to get task generator: ", taskType)
	return nil
}
