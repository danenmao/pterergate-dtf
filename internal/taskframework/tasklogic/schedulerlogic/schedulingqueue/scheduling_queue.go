package schedulingqueue

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskdef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/generationqueue"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/schedulerlogic/scheduler"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/tasklogic/tasklogicdef"
	"github.com/danenmao/pterergate-dtf/internal/tasktool"
)

const RRQueueIdx uint32 = 1000000
const PriorityBoostMaxTaskCount = uint32(10)

// 调度队列
type SchedulingQueue struct {
	QuotaGroupName string                    // 队列所属的资源组的名称
	QueueIndex     uint32                    // 队列在资源组内的索引
	QueueKeyName   string                    // 队列的Key名
	TimeSlice      uint32                    // 任务队列的轮转时间片，单位为ms
	BaseQueueSlice uint32                    // 任务队列授予任务的基础时间片的数量, 单位为个
	Scheduler      scheduler.IQueueScheduler // 调度队列的调度接口
	NextQueue      *SchedulingQueue          // 下一个调度队列
	TaskCount      uint                      // 队列中的任务数
}

// 创建一个优先级队列
func NewPriorityQueue(groupName string, queueName string, idx uint32) *SchedulingQueue {
	return &SchedulingQueue{
		QuotaGroupName: groupName,
		QueueIndex:     idx,
		QueueKeyName:   queueName,
		TimeSlice:      QueueBaseTimeSlice + idx*QueueTimeSliceStep,
		BaseQueueSlice: PriorityBaseQueueSliceCount,
		Scheduler:      &scheduler.FCFSScheduler{QueueKeyName: queueName},
	}
}

// 创建一个低优先级队列
func NewRRQueue(groupName string, queueName string) *SchedulingQueue {
	return &SchedulingQueue{
		QuotaGroupName: groupName,
		QueueIndex:     RRQueueIdx,
		QueueKeyName:   queueName,
		TimeSlice:      RRQueueTimeSlice,
		BaseQueueSlice: 1000000,
		Scheduler:      &scheduler.RRScheduler{QueueKeyName: queueName},
	}
}

// 调度任务
// 从当前的调度队列中选出一个任务, 进行调度, 返回获取的子任务列表
// 处理调度队列为空的情况, retTaskId为0, subtasks返回的元素为空
func (queue *SchedulingQueue) Schedule(
	retTaskId *taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
) (bool, error) {

	// 从队列中选择要调度的任务
	var taskId taskmodel.TaskIdType = 0
	noTask := false
	err := queue.Scheduler.Schedule(&taskId, &noTask)
	if err != nil {
		glog.Warning("failed to schedule a task: ", err.Error())
		return false, err
	}

	// 调度队列中没有任务, 返回空列表
	if noTask {
		*retTaskId = 0
		*subtasks = []taskmodel.SubtaskBody{}
		return false, nil
	}

	// 将调度的当前任务记录到当前任务列表中
	AddToCurrentTaskList(taskId)

	// 取任务的子任务列表
	finished := false
	quietTask := false
	err = queue.getSubtasks(taskId, subtasks, &finished, &quietTask)
	if err != nil {
		glog.Warning("failed to get subtasks of task: ", taskId, ", ", err.Error())
		//return false, err
	}

	// 如果生成完成，将任务从调度队列中移除
	if finished {
		glog.Info("task generation finished: ", taskId)
		queue.RemoveTask(taskId)
		RemoveFromCurrentTaskListDirectly(taskId)

		*retTaskId = taskId
		return false, nil
	}

	// 获取任务在当前队列的剩余时间片数量
	remainSliceCount, err := queue.getTaskRemainSliceCount(taskId)
	if err != nil {
		glog.Warning("failed to get remain slice count of task: ", taskId, err.Error())
		remainSliceCount = 2
	}

	glog.Infof("task %d remain time slice: %d", taskId, remainSliceCount)
	exhausted := false
	if remainSliceCount > 1 {
		// 若任务还有时间片, 将任务移到调度队列尾部
		queue.MoveTaskToTail(taskId, !quietTask)
	} else {
		// 若任务的时间片数量耗尽，返回状态, 由调用者处理
		exhausted = true
	}

	// 返回
	*retTaskId = taskId
	return exhausted, nil
}

// 各队列尾部添加任务
func (queue *SchedulingQueue) AppendTask(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	priority uint32,
) error {

	// 设置任务在本队列的调度数据
	err := queue.setTaskScheduleData(taskId, priority)
	if err != nil {
		glog.Warning("failed to set task schedule data: ", taskId, ",", err)
		return err
	}

	// 将任务添加到队列尾部
	pipeline := redistool.DefaultRedis().Pipeline()
	pipeline.RPush(context.Background(), queue.QueueKeyName, uint64(taskId))
	RemoveFromCurrentTaskList(taskId, pipeline)
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to append task to queue key: ", taskId, ", ", err)
		return err
	}

	queue.TaskCount += 1
	glog.Info("succeeded to append task to queue: ", taskId, ",", queue.QueueKeyName)
	return nil
}

// 设置任务在本队列的调度数据
func (queue *SchedulingQueue) setTaskScheduleData(
	taskId taskmodel.TaskIdType,
	priority uint32,
) error {

	// 任务是初次加入队列中,
	// 根据任务的优先级计算任务的初始时间片数量
	initSlice := queue.calcTaskSliceCount(priority)

	// 读取任务的调度数据
	data := tasklogicdef.TaskScheduleData{
		ResourceGroupName:   queue.QuotaGroupName,
		CurrentQueue:        queue.QueueIndex,
		CurrentQueueKeyName: queue.QueueKeyName,
		InitiallQueueSlice:  initSlice,
		QueueSlice:          initSlice,
	}

	// 保存任务的调度数据
	err := tasktool.SaveTaskScheduleData(taskId, &data)
	if err != nil {
		glog.Warning("failed to save task schedule data: ", taskId, ", ", err)
		return err
	}

	glog.Info("succeeded to set task schedule data on queue: ", taskId, ",", queue.QueueKeyName)
	return nil
}

// 根据优先级得到任务的时间片数量
func (queue *SchedulingQueue) calcTaskSliceCount(priority uint32) uint32 {

	priorityBonus := uint32(0)
	if priority <= taskdef.TaskPriority_Low {
		priorityBonus = PriorityBonus_Low
	} else if priority <= taskdef.TaskPriority_Medium {
		priorityBonus = PriorityBonus_Medium
	} else {
		priorityBonus = PriorityBonus_High
	}

	return priorityBonus * queue.BaseQueueSlice
}

// 从队列中移除任务
func (queue *SchedulingQueue) RemoveTask(taskId taskmodel.TaskIdType) error {

	queue.TaskCount -= 1
	glog.Info("succeeded to remove task from queue: ", taskId, ",", queue.QueueKeyName)
	return nil
}

// 向队列尾部添加PriorityBoost的任务列表
func (queue *SchedulingQueue) AppendBoostTaskList(taskIdList *[]taskmodel.TaskIdType) error {

	// 设置这些任务在本队列的调度数据
	for _, task := range *taskIdList {
		queue.setTaskScheduleData(task, taskdef.TaskPriority_Low)
	}

	// 批量添加到队列尾部
	vals := []interface{}{}
	for _, task := range *taskIdList {
		vals = append(vals, task)
	}

	pipeline := redistool.DefaultRedis().Pipeline()
	pipeline.RPush(context.Background(), queue.QueueKeyName, vals...)
	RemoveListFromCurrentTaskList(vals, pipeline)
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to push boost task list: ", err)
		return err
	}

	glog.Info("succeeded to append boost taks list: ", queue.QueueKeyName, *taskIdList)
	return nil
}

// 从队列中取出前若干个任务
func (queue *SchedulingQueue) PopBoostTask(taskIdList *[]taskmodel.TaskIdType) error {

	// 构造命令pipeline
	pipeline := redistool.DefaultRedis().Pipeline()
	for i := uint32(0); i < PriorityBoostMaxTaskCount; i++ {
		pipeline.LPop(context.Background(), queue.QueueKeyName)
	}

	cmdList, err := pipeline.Exec(context.Background())
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		glog.Warning("failed to exec pop boost task list pipeline: ", queue.QueueKeyName,
			",", err)
		return err
	}

	// 从pipeline结果中读取任务ID列表
	for _, cmd := range cmdList {
		err = cmd.Err()

		// 无元素，退出
		if err == redis.Nil {
			break
		}

		if err != nil {
			glog.Warning("pipeline cmd return error: ", err)
			break
		}

		strCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			glog.Error("invalid cmd type: ", cmd)
			continue
		}

		taskId, err := strCmd.Uint64()
		if err != nil {
			glog.Warning("failed to convert task id: ", strCmd.Val(), ",", err)
			continue
		}

		*taskIdList = append(*taskIdList, taskmodel.TaskIdType(taskId))
	}

	return nil
}

// 将任务移到队尾
func (queue *SchedulingQueue) MoveTaskToTail(taskId taskmodel.TaskIdType, toDecrSlice bool) error {

	// 减少任务的时间片数量
	if toDecrSlice {
		err := queue.DecreaseTaskSliceCount(taskId)
		if err != nil {
			glog.Warning("failed to decrease task slice count: ", taskId, ",", err)
		}
	}

	// 移到的尾部
	pipeline := redistool.DefaultRedis().Pipeline()
	pipeline.RPush(context.Background(), queue.QueueKeyName, uint64(taskId))
	RemoveFromCurrentTaskList(taskId, pipeline)
	_, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec move task to queue tail pipeline: ", taskId, ",", err)
		return err
	}

	glog.Info("succeeded to move task to queue tail: ", taskId, ",", queue.QueueKeyName)
	return nil
}

// 减少任务的时间片数量
func (queue *SchedulingQueue) DecreaseTaskSliceCount(
	taskId taskmodel.TaskIdType,
) error {

	// 取当前的调度数据
	data := tasklogicdef.TaskScheduleData{}
	err := tasktool.GetTaskScheduleData(taskId, &data)
	if err != nil {
		glog.Warning("failed to get task schedule data when move task to tail: ", taskId, ",", err)
		return err
	}

	// 减少本次花费的时间片数量
	data.QueueSlice -= 1

	// 保存调度数据
	err = tasktool.SaveTaskScheduleData(taskId, &data)
	if err != nil {
		glog.Warning("failed to update task schedule data: ", taskId, ",", err)
		return err
	}

	glog.Info("succeeded to decrease task slice count: ", taskId)
	return nil
}

// 获取任务在队列中剩余的时间片数量
func (queue *SchedulingQueue) getTaskRemainSliceCount(
	taskId taskmodel.TaskIdType,
) (uint32, error) {

	// 取任务的调度数据
	var scheduleData = tasklogicdef.TaskScheduleData{}
	err := tasktool.GetTaskScheduleData(taskId, &scheduleData)
	if err != nil {
		glog.Warning("failed to get task schedule data: ", taskId, ", ", err.Error())
		return 0, err
	}

	// 返回任务在队列中的剩余时间片数量
	glog.Info("task schedule data: ", scheduleData)
	return scheduleData.QueueSlice, nil
}

// 获取任务的子任务列表
func (queue *SchedulingQueue) getSubtasks(
	taskId taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
	retFinished *bool,
	retQuietTask *bool,
) error {

	// 循环取子任务
	err := queue.pickSubtaskLoop(taskId, subtasks, retFinished, retQuietTask)
	if err != nil {
		glog.Warning("failed to loop subtasks: ", taskId, ", ", err.Error())
		return err
	}

	return nil
}

// 取子任务循环
func (queue *SchedulingQueue) pickSubtaskLoop(
	taskId taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
	retFinished *bool,
	retQuietTask *bool,
) error {

	// 初始化记录
	start := time.Now()
	subtaskCount := 0

	// 取任务的调度数据
	var scheduleData = tasklogicdef.TaskScheduleData{}
	err := tasktool.GetTaskScheduleData(taskId, &scheduleData)
	if err != nil {
		glog.Warning("failed to get task schedule data: ", taskId, ", ", err.Error())
	}

	// 取子任务循环
	for {
		// 检查消耗的时间片
		now := time.Now()
		if uint32(now.Sub(start).Milliseconds()) >= queue.TimeSlice {
			glog.Info("task time slice is exhausted: ", taskId)
			break
		}

		// 取子任务
		finished := false
		var subtaskData = taskmodel.SubtaskBody{}
		err := GetSubtask(taskId, &subtaskData, &finished)
		if err != nil && err != errordef.ErrNotFound {
			glog.Warning("failed to get subtask: ", taskId, ", ", err.Error())
			continue
		}

		gotSubtask := (err == nil)
		if gotSubtask {
			// 为子任务创建key
			tasktool.CreateSubtaskInfoKey(uint64(subtaskData.SubtaskId), &subtaskData)

			// 将子任务添加到任务的子任务列表中
			tasktool.AddSubtaskToTask(taskId, uint64(subtaskData.SubtaskId))
			*subtasks = append(*subtasks, subtaskData)
			subtaskCount++
		}

		// 处理生成完成的情况
		if finished {
			*retFinished = true
			glog.Info("task generation finished: ", taskId)
			break
		}

		// 处理任务静默期的情况
		if !gotSubtask && subtaskCount == 0 {
			// 初次触发静默期，更新调度记录
			if scheduleData.QuietStartTime == 0 {
				scheduleData.QuietStartTime = uint64(now.Unix())
				tasktool.SaveTaskScheduleData(taskId, &scheduleData)
			}

			// 若未超出静默期限制，任务处于静默期内
			if uint64(now.Unix())-scheduleData.QuietStartTime < QuietTaskMaxInterval {
				*retQuietTask = true
				glog.Info("quiet task found: ", taskId)
				break
			} else {
				glog.Info("exceeds max quiet task interval: ", taskId)
				break
			}
		}

	} // for

	return nil
}

// get a subtask from the subtask queue
func GetSubtask(
	taskId taskmodel.TaskIdType,
	subtaskData *taskmodel.SubtaskBody,
	finished *bool,
) error {

	// get a subtask from the subtask queue
	queue := generationqueue.GenerationQueue{TaskId: taskId}
	err := queue.Pop(subtaskData)
	noSubtask := (err == errordef.ErrNotFound)

	// check if the generation is over
	*finished = tasktool.CheckIfTaskGenerationCompleted(taskId)
	if *finished {
		subtaskCount, err := queue.GetSubtaskCount(taskId)
		if err == nil && subtaskCount > 0 {
			*finished = false
		}
	}

	// no subtask
	if noSubtask {
		return errordef.ErrNotFound
	}

	if err != nil {
		glog.Warning("failed to pop subtask: ", taskId, ",", err)
		return err
	}

	return nil
}
