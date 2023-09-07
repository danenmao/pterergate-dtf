package generator

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/taskframework/taskflow/flowdef"
	"pterergate-dtf/internal/tasktool"
)

// 协程, 检查并处理要生成的任务，执行生成操作
func StartGenerateTaskRoutine() {

	// 检查当前实例生成的任务数是否超过上限
	if CheckIfExceedLimit() {
		glog.Info("exceed task generation limit")
		return
	}

	// 取要生成的任务id
	taskId, err := getTaskIdToGenerate()
	if err == errordef.ErrNotFound {
		return
	}

	if err != nil {
		glog.Warning("failed to get task id to schedule: ", err)
		return
	}

	if taskId == 0 {
		glog.Error("zero task id")
		return
	}

	// 对于取到的 $taskid, 创建协程执行生成操作.
	startGeneration(taskId)
}

// 获取要执行生成的任务ID
func getTaskIdToGenerate() (taskmodel.TaskIdType, error) {

	// 按优先级从高到低从待生成任务列表中取要执行调度的任务ID
	opt := redis.ZRangeBy{
		Min: "-inf", Max: "+inf",
		Offset: 0, Count: 1,
	}

	cmd := redistool.DefaultRedis().ZRangeByScore(
		context.Background(), config.ToGenerateTaskZset, &opt,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get a task to generate: ", err)
		return 0, err
	}

	// 如果列表为空，表示没有待生成的任务
	taskList := cmd.Val()
	if len(taskList) == 0 {
		return 0, errordef.ErrNotFound
	}

	glog.Info("got a task to generate: ", taskList[0])

	// 尝试获取生成的计数
	if !IncrIfNotExceedLimit() {
		glog.Info("exceed the limit, cannot get a generation routine: ", taskList[0])
		return 0, errordef.ErrNotFound
	}

	toDecrGeneratingCount := true
	defer func() {
		if toDecrGeneratingCount {
			DecrGeneratingRoutineCount()
		}
	}()

	// 尝试删除，如果删除成功，则获取了此元素
	remCmd := redistool.DefaultRedis().ZRem(context.Background(), config.ToGenerateTaskZset, taskList[0])
	err = remCmd.Err()
	if err != nil {
		glog.Warning("failed to rem task from to-generate list: ", taskList[0])
		return 0, err
	}

	// 如果返回1，表示删除成功，获取了些任务ID; 返回0, 表示元素不存在，被其他实例删除
	if remCmd.Val() == 0 {
		glog.Info("task removed by other routine: ", taskList[0])
		return 0, errordef.ErrNotFound
	}

	// 取返回的task id
	taskIdStr := taskList[0]
	var taskId uint64 = 0
	taskId, err = strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		glog.Warning("failed to convert task id: ", taskIdStr)
		return 0, err
	}

	if taskId == 0 {
		glog.Warning("zero task id from to-generate queue")
		return 0, errordef.ErrNotFound
	}

	toDecrGeneratingCount = false
	glog.Info("got a task to generate: ", taskId)
	return taskmodel.TaskIdType(taskId), nil
}

// 开始生成流程
func startGeneration(taskId taskmodel.TaskIdType) {

	err := TryToOwnTask(taskId)
	if err != nil {
		glog.Warning("failed to own task: ", taskId, err)
	}

	// 启动生成例程
	glog.Info("to start a generation routine: ", taskId)
	go TaskGenerationRoutine(taskId, false)
}

// 任务的生成例程
func TaskGenerationRoutine(taskId taskmodel.TaskIdType, toRecover bool) {

	glog.Info("begin to generate task: ", taskId, ", ", toRecover)

	// 减少正在生成的例程数
	defer DecrGeneratingRoutineCount()

	// 释放对任务的所有权
	defer ReleaseTask(taskId)

	// 判断任务类型
	taskType := uint32(0)
	err := tasktool.GetTaskType(taskmodel.TaskIdType(taskId), &taskType)
	if err != nil {
		glog.Warning("failed to get task type, return: ", taskId, ", ", err)
		return
	}

	// 根据任务类型，执行不同的生成逻辑
	glog.Info("task being generated type: ", taskId, ", ", taskType)
	taskGenerationLogic(taskId, toRecover, taskType)

	// 清理操作
	glog.Info("finished to generate task: ", taskId)
}

func taskGenerationLogic(taskId taskmodel.TaskIdType, toRecover bool, taskType uint32) {

	glog.Info("begin to generate a plugin task: ", taskId, ", ", taskType)

	// 获取任务的信息，信息在创建任务的流程中提供
	createParam := flowdef.TaskCreateParam{}
	err := tasktool.GetTaskCreateParam(taskmodel.TaskIdType(taskId), &createParam)
	if err != nil {
		glog.Warning("failed to get image task create param: ", taskId, ", ", err)
		return
	}

	if !toRecover {
		// 将任务添加到调度队列中
	}

	// 执行生成逻辑
	step := uint32(0)
	err = InitGeneration(taskId, &step)
	if err == nil {
		GenerationLoop(taskId, &createParam)
	}

	FinishGeneration(taskId)

}

// 任务插件的生成逻辑
func GenerationLoop(
	taskId taskmodel.TaskIdType,
	createParam *flowdef.TaskCreateParam,
) {
}

func InitGeneration(taskId taskmodel.TaskIdType, step *uint32) error {
	return nil
}

// 完成任务生成操作
func FinishGeneration(taskId taskmodel.TaskIdType) error {
	return nil
}
