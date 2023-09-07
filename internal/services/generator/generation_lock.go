package generator

import (
	"time"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/tasktool"
)

// 尝试获取对任务生成的所有权
func TryToOwnTask(taskId taskmodel.TaskIdType) error {

	err := redistool.GetLockWithExpire(tasktool.GetTaskLockKey(taskId), 200, time.Minute*5)
	return err
}

// 释放对任务生成的所有权
func ReleaseTask(taskId taskmodel.TaskIdType) {
	redistool.ReleaseLock(tasktool.GetTaskLockKey(taskId))
}

// 对所有权续期
func RenewTask(taskId taskmodel.TaskIdType) {
	redistool.RenewLock(tasktool.GetTaskLockKey(taskId), time.Second*60)
}
