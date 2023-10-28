package generator

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/generationlogic"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/schedulerlogic"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/tasklogicdef"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

// 协程, 检查并处理要生成的任务，执行生成操作
func StartTaskGenerationRoutine() {
	// 检查当前实例生成的任务数是否超过上限
	if IsFull() {
		glog.Warning("exceed task generation limit")
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

	// 对于取到的taskid, 创建协程执行生成操作.
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
	if !IncrIfNotFull() {
		glog.Info("exceed the limit, cannot get a generation routine: ", taskList[0])
		return 0, errordef.ErrNotFound
	}

	toDecrGeneratingCount := true
	defer func() {
		if toDecrGeneratingCount {
			Decr()
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
	err := tasktool.TryToOwnTask(taskId)
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
	defer Decr()

	// 释放对任务的所有权
	defer tasktool.ReleaseTask(taskId)

	// 判断任务类型
	taskType := uint32(0)
	err := tasktool.GetTaskType(taskmodel.TaskIdType(taskId), &taskType)
	if err != nil {
		glog.Warning("failed to get task type, return: ", taskId, ", ", err)
		return
	}

	// 根据任务类型，执行不同的生成逻辑
	glog.Info("task being generated type: ", taskId, ", ", taskType)
	taskGenerationImpl(taskId, toRecover, taskType)

	// 清理操作
	glog.Info("finished to generate task: ", taskId)
}

func taskGenerationImpl(taskId taskmodel.TaskIdType, toRecover bool, taskType uint32) {

	glog.Info("begin to generate a plugin task: ", taskId, ", ", taskType)

	// 获取任务的信息，信息在创建任务的流程中提供
	createParam := tasklogicdef.TaskCreateParam{}
	err := tasktool.GetTaskCreateParam(taskmodel.TaskIdType(taskId), &createParam)
	if err != nil {
		glog.Warning("failed to get image task create param: ", taskId, ", ", err)
		return
	}

	if !toRecover {
		// 将任务添加到调度队列中
		err = AddTaskToScheduler(taskId, createParam.ResourceGroupName, createParam.TaskType,
			createParam.Priority)
		if err != nil {
			glog.Warning("failed to add image task to scheduler: ", taskId, ", ", err)
			return
		}
	}

	// 执行生成逻辑
	step := uint32(0)
	err = InitGeneration(taskId, &step)
	if err == nil {
		GenerationMainLoop(taskId, &createParam)
	}

	FinishGeneration(taskId)
}

// 将文件添加到调度队列中
func AddTaskToScheduler(
	taskId taskmodel.TaskIdType,
	groupName string,
	taskType uint32,
	priority uint32,
) error {
	return schedulerlogic.AddTaskToScheduler(taskId, groupName, taskType, priority)
}

// 任务插件的生成逻辑
func GenerationMainLoop(
	taskId taskmodel.TaskIdType,
	createParam *tasklogicdef.TaskCreateParam,
) {

	// 创建生成工作流
	flow := generationlogic.NewGenerationiFlow()

	// 初始化生成操作
	taskData := taskmodel.TaskParam{
		Creator:       taskmodel.TaskCreator{},
		ResourceGroup: createParam.ResourceGroupName,
		Priority:      createParam.Priority,
		TaskType:      createParam.TaskType,
		Timeout:       time.Duration(createParam.Timeout) * time.Second,
		TypeParam:     createParam.TypeParam,
	}

	err := flow.InitGeneration(taskmodel.TaskIdType(taskId), createParam.TaskType, &taskData)
	if err != nil {
		glog.Warning("failed to init task generation: ", taskId, createParam.TaskType, ",", err)
		return
	}

	// 执行生成循环
	err = flow.GenerationLoop()
	if err != nil {
		glog.Warning("task generation loop failed: ", taskId, ", ", err)
	}

	// 结束生成操作
	err = flow.FinishGeneration()
	if err != nil {
		glog.Warning("failed to finish task generation: ", taskId, ", ", err)
	}
}

func InitGeneration(taskId taskmodel.TaskIdType, step *uint32) error {

	if step == nil {
		panic("invalid step pointer")
	}

	// 检查 next_check_time值是否存在
	var currentStep uint32 = 0
	var toGenerate bool = false
	err := CheckGenerationStatus(taskId, &toGenerate, &currentStep)
	if err != nil {
		glog.Warning("failed to check generation status: ", taskId, err)
		return err
	}

	// 有其他生成协程在处理, 退出
	if !toGenerate {
		glog.Info("task be generating by other")
		return errors.New("task be generating by other")
	}

	// 返回任务之前生成逻辑的进展
	*step = currentStep
	glog.Info("former task generation step: ", taskId, currentStep)

	pipeline := redistool.DefaultRedis().Pipeline()

	// 创建 redis_task_generation.$taskid.progress, 更新next_check_time
	progressMap := map[string]interface{}{
		config.TaskGenerationKey_NextCheckTimeField: uint64(time.Now().Add(time.Minute).Unix()),
		config.TaskGenerationKey_StepField:          currentStep,
	}
	pipeline.HMSet(context.Background(), tasktool.GetTaskGenerationProgressKey(taskId), progressMap)

	// 将 $taskid 移入 redis_task_generation_zset，按照插入时间排序，表示任务进入了生成状态。
	pipeline.ZAdd(context.Background(), config.GeneratingTaskZset, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: taskId,
	})

	// 执行pipeline
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", taskId, err)
		return err
	}

	glog.Info("succeeded to init task generation: ", taskId, currentStep)
	return nil
}

// 完成任务生成操作
func FinishGeneration(taskId taskmodel.TaskIdType) error {

	pipeline := redistool.DefaultRedis().Pipeline()

	// 从 redis_task_generation_zset 中移除 $taskid
	pipeline.ZRem(context.Background(), config.GeneratingTaskZset, taskId)

	// 推入 redis_task_schedule_zset中
	pipeline.ZAdd(context.Background(), config.RunningTaskZset, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: taskId,
	})

	// 设置redis_task_generation.$taskid.progress 12小时后过期
	pipeline.Expire(context.Background(), tasktool.GetTaskGenerationProgressKey(taskId), time.Hour*12)

	// 设置任务生成完成的标记. redis_task_info.$taskinfo,
	// task_generation_completed = 1.
	pipeline.HSet(context.Background(), tasktool.GetTaskInfoKey(taskId), config.TaskInfo_GenerationCompletedField, 1)

	// 执行
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	glog.Info("succeeded to generate task: ", taskId)
	return nil
}

// 检查任务生成的状态
func CheckGenerationStatus(taskId taskmodel.TaskIdType, toGenerate *bool, currentStep *uint32) error {

	// 读取redis_task_generation.$taskid.progress
	cmd := redistool.DefaultRedis().HGetAll(context.Background(), tasktool.GetTaskGenerationProgressKey(taskId))
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get task generation progress key: ", taskId, err)
		return err
	}

	// map为空，表示key不存在，可以执行生成流程
	valMap := cmd.Val()
	if len(valMap) == 0 {
		glog.Info("empty task generation progress")
		*toGenerate = true
		*currentStep = 0
		return nil
	}

	// 检查 next_check_time值是否存在,或已过期
	nextCheckTimeStr, ok := valMap[config.TaskGenerationKey_NextCheckTimeField]
	if !ok {
		glog.Info("no next_check_time field")
		*toGenerate = true
		*currentStep = 0
		return nil
	}

	nextCheckTime, err := strconv.ParseUint(nextCheckTimeStr, 10, 64)
	if err != nil {
		glog.Warning("failed to convert next_check_time: ", nextCheckTimeStr, err)
		return err
	}

	// 若存在且未过期,表示有其他协程在处理, 退出
	if nextCheckTime >= uint64(time.Now().Unix()) {
		glog.Info("task generation not expired: ", taskId)
		*toGenerate = false
		return nil
	}

	// 已过期，需要进行处理，获取step
	*toGenerate = true
	stepStr, ok := valMap[config.TaskGenerationKey_StepField]
	if !ok {
		glog.Info("no step field")
		*currentStep = 0
		return nil
	}

	step, err := strconv.Atoi(stepStr)
	if err != nil {
		glog.Warning("failed to convert step: ", stepStr, err)
		return err
	}

	*currentStep = uint32(step)
	glog.Info("found task generation step: ", taskId, step)
	return nil
}

// 更新生成的step值
func RefreshTaskGenerationStep(taskId taskmodel.TaskIdType, step uint32) error {

	cmd := redistool.DefaultRedis().HSet(
		context.Background(),
		tasktool.GetTaskGenerationProgressKey(taskId),
		config.TaskGenerationKey_StepField,
		step,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to refresh task generation step value: ", taskId, step, err)
		return err
	}

	return nil
}
