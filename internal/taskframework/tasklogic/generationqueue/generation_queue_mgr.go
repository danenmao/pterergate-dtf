package generationqueue

import (
	"errors"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

// 管理当前生成器实例下的所有任务的子任务队列
type GenerationiQueueMgr struct {
	GenerationQueueMap map[taskmodel.TaskIdType]*GenerationQueue // 各任务的子任务队列表
}

// 添加任务
// 当创建任务、恢复任务生成时，执行添加操作
func (mgr *GenerationiQueueMgr) AddTask(taskId taskmodel.TaskIdType) error {
	if taskId == 0 {
		glog.Info("invalid task id: ", taskId)
		return errors.New("invalid task id")
	}

	mgr.GenerationQueueMap[taskId] = &GenerationQueue{
		TaskId: taskId,
	}

	return nil
}

// 删除任务
// 当任务完成时，执行删除操作
func (mgr *GenerationiQueueMgr) RemoveTask(taskId taskmodel.TaskIdType) error {
	// 检查任务是否存在
	_, ok := mgr.GenerationQueueMap[taskId]
	if !ok {
		glog.Warning("no task id found: ", taskId)
		return nil
	}

	// 删除记录
	delete(mgr.GenerationQueueMap, taskId)
	glog.Info("succeeded to remove task subtask queue: ", taskId)

	return nil
}

// 将子任务放入子任务队列中
func (mgr *GenerationiQueueMgr) PushSubtask(
	taskId taskmodel.TaskIdType,
	subtask *taskmodel.SubtaskBody,
) error {
	queue, ok := mgr.GenerationQueueMap[taskId]
	if !ok {
		glog.Warning("task id not found in subtask queue map: ", taskId)
		return errors.New("task id not found in subtask queue map")
	}

	// 将子任务放到任务的子任务队列中
	queue.Push(subtask)

	return nil
}

// 从子任务队列中取子任务
func (mgr *GenerationiQueueMgr) PopSubtask(
	taskId taskmodel.TaskIdType,
	subtask *taskmodel.SubtaskBody,
) error {
	queue, ok := mgr.GenerationQueueMap[taskId]
	if !ok {
		glog.Warning("task id not found in subtask queue map: ", taskId)
		return errors.New("task id not found in subtask queue map")
	}

	// 从子任务队列中取子任务
	err := queue.Pop(subtask)
	if err == errordef.ErrNotFound {
		return err
	}

	return nil
}
