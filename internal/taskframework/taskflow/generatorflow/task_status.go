package generatorflow

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/redistool"
)

// 保存任务创建
const RedisTaskStatusKeyPrefix = "task.status."

// 获取保存任务调度数据的Key
func GetTaskStatusKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", RedisTaskStatusKeyPrefix, taskId)
}

// 根据任务ID从Redis中加载任务的生成状态
func LoadStatus(
	taskId taskmodel.TaskIdType,
	retTaskStatus *string,
) error {

	keyName := GetTaskStatusKey(taskId)
	cmd := redistool.DefaultRedis().Get(context.Background(), keyName)
	err := cmd.Err()

	// 要处理key不存在的场景。
	// key不存在，表示任务未执行过
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		glog.Warning("failed to get task status key: ", taskId, ", ", err.Error())
		return err
	}

	*retTaskStatus = cmd.Val()
	return nil
}

// 保存任务的运行状态数据
func SaveStatus(
	taskId taskmodel.TaskIdType,
	taskStatus string,
) error {

	cmd := redistool.DefaultRedis().Set(
		context.Background(),
		GetTaskStatusKey(taskId),
		taskStatus,
		time.Hour*48,
	)

	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to save task status key: ", taskId, ", ", err.Error())
		return err
	}

	glog.Info("succeeded to save task status: ", taskId)
	return nil
}
