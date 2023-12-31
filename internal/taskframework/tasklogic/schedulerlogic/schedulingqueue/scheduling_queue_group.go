package schedulingqueue

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/misc"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/tasklogicdef"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

// 调度队列组
type SchedulingTeam struct {
	TeamName       string             // 调度队列组名
	PriorityQueues []*SchedulingQueue // 优先级队列, 队列内的任务有优先级
	RRQueue        *SchedulingQueue   // 低优先级队列, 队列内的任务使用时间片轮转策略
}

// 初始化
func (queues *SchedulingTeam) Init(groupName string) error {
	// 记录队列组名
	queues.TeamName = groupName

	// 创建调度队列组
	err := queues.createSchedulingQueues()
	if err != nil {
		glog.Warning("failed to create scheduling queues: ", err.Error())
		return err
	}

	// 启动每个优先级队列的工作例程
	for i := 1; i < len(queues.PriorityQueues); i++ {
		go queues.priorityBoostRoutine(uint32(i))
	}

	// 启用RR队列的工作例程
	go queues.rrPriorityBoost()
	go queues.remainAcceleration()

	glog.Infof("succeeded to init scheduling queue array: %+v", queues)
	misc.DumpDataInTest("scheduling queue array", queues)
	return nil
}

// 获取调度队列组中的任务数
func (queues *SchedulingTeam) GetTaskCount() (taskCount uint, err error) {
	taskCount = 0
	for _, queue := range queues.PriorityQueues {
		taskCount += queue.TaskCount
	}

	taskCount += queues.RRQueue.TaskCount
	return
}

// 向调度队列组中添加任务
func (queues *SchedulingTeam) AddTask(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	priority uint32,
) error {

	// 任务的调度数据
	// 将任务添加到高优先级队列中
	err := queues.PriorityQueues[0].AppendTask(taskId, taskType, priority)
	if err != nil {
		glog.Warning("failed to append task to the highest priority queue: ", taskId, err.Error())
		return err
	}

	glog.Info("succeeded to append task to the highest priority queue: ", taskId)
	return nil
}

// 从调度队列组中选出一个任务来执行, 返回任务的子任务列表
// 处理调度队列为空的情况, retTaskId为0, subtasks返回的元素为空
func (queues *SchedulingTeam) Schedule(
	retTaskId *taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
) error {
	var taskId taskmodel.TaskIdType = 0

	// 按优先级从各个调度队列中选出一个任务
	for _, queue := range queues.PriorityQueues {
		// 从当前的调度队列中选出一个任务, 进行调度
		exhausted, err := queue.Schedule(&taskId, subtasks)
		if err != nil {
			glog.Warning("failed to schedule a task in: ", queue.QueueKeyName, ", ", err.Error())
			continue
		}

		// 若任务的时间片数量耗尽，移至下一个队列中
		if exhausted {
			glog.Info("time slice of task is over: ", taskId, ",", queue.QueueKeyName)
			queue.RemoveTask(taskId)
			queues.appendToNextQueue(queue, taskId)
		}

		// 已选出任务, 返回子任务列表
		if len(*subtasks) > 0 {
			*retTaskId = taskId
			glog.Info("succeeded to schedule a task in: ", queue.QueueKeyName, ", ", taskId, len(*subtasks))
			return nil
		}
	} // for

	// 无法选出任务，执行低优先队列中的任务
	_, err := queues.RRQueue.Schedule(&taskId, subtasks)
	if err != nil {
		glog.Warning("failed to schedule a task in RR queue: ", err.Error())
		return err
	}

	*retTaskId = taskId
	if taskId != 0 {
		glog.Info("succeeded to schedule a task in RR queue, ", taskId, len(*subtasks))
	}

	return nil
}

// 将任务移到下个队列
func (queues *SchedulingTeam) appendToNextQueue(
	queue *SchedulingQueue,
	taskId taskmodel.TaskIdType,
) error {
	if queue.NextQueue == nil {
		glog.Warning("null next queue pointer")
		return nil
	}

	// 取任务的类型信息
	createParam := tasklogicdef.TaskCreateParam{}
	err := tasktool.GetTaskCreateParam(taskId, &createParam)
	if err != nil {
		glog.Warning("failed to get task create param: ", taskId, ", ", err)
		return err
	}

	// 将任务添加到下个队列的尾部
	err = queue.NextQueue.AppendTask(taskId, createParam.TaskType, createParam.Priority)
	if err != nil {
		glog.Warning("failed to transfer a task to next queue: ", taskId, err.Error())
		return err
	}

	glog.Info("succeeded to append task to next queue: ", taskId, ", current:", queue.QueueKeyName,
		",next:", queue.NextQueue.QueueKeyName)
	return nil
}

// 创建调度队列
func (queues *SchedulingTeam) createSchedulingQueues() error {

	// 创建优先级调度队列
	for i := uint32(0); i < PrioirtyQueueCount; i++ {
		queueName := fmt.Sprintf("%s.%s.P%d.queue", ScheduleQueueKeyPrefix, queues.TeamName, i)
		queues.PriorityQueues = append(queues.PriorityQueues, NewPriorityQueue(queues.TeamName, queueName, i))
	}

	// 创建低优先级调度队列
	queueName := fmt.Sprintf("%s.%s.RR.queue", ScheduleQueueKeyPrefix, queues.TeamName)
	queues.RRQueue = NewRRQueue(queues.TeamName, queueName)

	// 设置队列的next queue
	for i := uint32(0); i < PrioirtyQueueCount-1; i++ {
		queues.PriorityQueues[i].NextQueue = queues.PriorityQueues[i+1]
	}

	queues.PriorityQueues[PrioirtyQueueCount-1].NextQueue = queues.RRQueue
	queues.RRQueue.NextQueue = nil

	glog.Info("succeeded to create a scheduling group")
	return nil
}

// Priority Boost策略例程
func (queues *SchedulingTeam) priorityBoostRoutine(idx uint32) error {
	routine.ExecRoutineWithInterval(
		"priorityBoostRoutine",
		func() {
			queues.triggerPriorityBoost(idx)
		},
		time.Duration(PriorityBoostInterval)*time.Second,
	)

	return nil
}

// RR队列的Priority Boost策略例程
func (queues *SchedulingTeam) rrPriorityBoost() error {
	routine.ExecRoutineWithInterval(
		"rrPriorityBoostRoutine",
		func() {
			queues.triggerRRPriorityBoost()
		},
		time.Duration(RRPriorityBoostInterval)*time.Second,
	)

	return nil
}

// 任务剩余时间加速策略例程
func (queues *SchedulingTeam) remainAcceleration() error {
	routine.ExecRoutineWithInterval(
		"remainAccelerationRoutine",
		func() {
			queues.triggerRemainAcceleration()
		},
		time.Duration(RemainTaskAccelerationInteral)*time.Second,
	)

	return nil
}

// 执行Priority Boost策略
func (queues *SchedulingTeam) triggerPriorityBoost(idx uint32) error {
	if idx >= uint32(len(queues.PriorityQueues)) {
		glog.Error("priority queue idx out of range: ", idx)
		return nil
	}

	// 在队列上执行Priority Boost
	currentQueue := queues.PriorityQueues[idx]
	err := queues.priorityBoostOnQueue(currentQueue)
	if err != nil {
		glog.Warning("failed to exec priority boost on queue: ",
			currentQueue.QueueKeyName, ",", err)
	}

	glog.Info("succeeded to exec priority boost on queue: ", idx, ",", currentQueue.QueueKeyName)
	return nil
}

// 在指定的队列上执行Priority Boost策略
func (queues *SchedulingTeam) priorityBoostOnQueue(
	currentQueue *SchedulingQueue,
) error {
	queueName := currentQueue.QueueKeyName
	glog.Info("execute priority boost for queue: ", queueName)

	// 从队列中取出前若干个任务
	taskIdList := []taskmodel.TaskIdType{}
	err := currentQueue.PopBoostTask(&taskIdList)
	if err != nil {
		glog.Warning("failed to pop task list: ", queueName, ",", err)
		return err
	}

	if len(taskIdList) <= 0 {
		glog.Info("no task to boost on queue: ", queueName)
		return nil
	}

	// 添加到高优先级队列中
	AddListToCurrentTaskList(taskIdList)
	err = queues.PriorityQueues[0].AppendBoostTaskList(&taskIdList)
	if err != nil {
		glog.Warning("failed to append boost task list: ", queueName, ",", err)
	}

	return nil
}

// 对RR队列执行Priority Boost策略
func (queues *SchedulingTeam) triggerRRPriorityBoost() error {
	currentQueue := queues.RRQueue
	err := queues.priorityBoostOnQueue(currentQueue)
	if err != nil {
		glog.Warning("failed to exec priority boost on RR queue: ",
			currentQueue.QueueKeyName, ",", err)
	}

	glog.Info("succeeded to exec priority boost on RR queue: ", currentQueue.QueueKeyName)
	return nil
}

// 执行任务剩余时间加速策略
func (queues *SchedulingTeam) triggerRemainAcceleration() error {
	return nil
}
