package subtaskqueue

import (
	"errors"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
)

// 管理当前生成器实例下的所有任务的子任务队列
type SubtaskQueueMgr struct {
	SubtaskQueueMap map[taskmodel.TaskIdType]*SubtaskQueue // 各任务的子任务队列表
}

// 添加任务
// 当创建任务、恢复任务生成时，执行添加操作
func (mgr *SubtaskQueueMgr) AddTask(taskId taskmodel.TaskIdType) error {

	if taskId == 0 {
		glog.Info("invalid task id: ", taskId)
		return errors.New("invalid task id")
	}

	mgr.SubtaskQueueMap[taskId] = &SubtaskQueue{
		TaskId: taskId,
	}

	return nil
}

// 删除任务
// 当任务完成时，执行删除操作
func (mgr *SubtaskQueueMgr) RemoveTask(taskId taskmodel.TaskIdType) error {

	// 检查任务是否存在
	_, ok := mgr.SubtaskQueueMap[taskId]
	if !ok {
		glog.Warning("no task id found: ", taskId)
		return nil
	}

	// 删除记录
	delete(mgr.SubtaskQueueMap, taskId)
	glog.Info("succeeded to remove task subtask queue: ", taskId)

	return nil
}

// 将子任务放入子任务队列中
func (mgr *SubtaskQueueMgr) PushSubtask(
	taskId taskmodel.TaskIdType,
	subtask *taskmodel.SubtaskData,
) error {

	queue, ok := mgr.SubtaskQueueMap[taskId]
	if !ok {
		glog.Warning("task id not found in subtask queue map: ", taskId)
		return errors.New("task id not found in subtask queue map")
	}

	// 将子任务放到任务的子任务队列中
	queue.PushSubtask(subtask)

	return nil
}

// 从子任务队列中取子任务
func (mgr *SubtaskQueueMgr) PopSubtask(
	taskId taskmodel.TaskIdType,
	subtask *taskmodel.SubtaskData,
) error {

	queue, ok := mgr.SubtaskQueueMap[taskId]
	if !ok {
		glog.Warning("task id not found in subtask queue map: ", taskId)
		return errors.New("task id not found in subtask queue map")
	}

	// 从子任务队列中取子任务
	err := queue.PopSubtask(subtask)
	if err == errordef.ErrNotFound {
		return err
	}

	return nil
}
