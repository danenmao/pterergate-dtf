package tasktool

import (
	"time"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

// 尝试获取对任务生成的所有权
func TryToOwnTask(taskId taskmodel.TaskIdType) error {
	err := redistool.LockWithExpire(GetTaskLockKey(taskId), 200, time.Minute*5)
	return err
}

// 释放对任务生成的所有权
func ReleaseTask(taskId taskmodel.TaskIdType) {
	redistool.Unlock(GetTaskLockKey(taskId))
}

// 对所有权续期
func RenewTask(taskId taskmodel.TaskIdType) {
	redistool.RenewLock(GetTaskLockKey(taskId), time.Second*60)
}
