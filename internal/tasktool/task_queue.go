package tasktool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/redistool"
)

// 将任务添加到创建队列
// 将 $taskid 推入 redis_creating_task_zset
func AddTaskToCreatingQueue(taskId taskmodel.TaskIdType) error {

	// 将 $taskid 推入 redis_creating_task_zset
	var z = redis.Z{
		Score:  float64(time.Now().Add(time.Second * time.Duration(config.EnvTaskCreatingTimeout)).Unix()),
		Member: taskId,
	}

	cmd := redistool.DefaultRedis().ZAdd(context.Background(), config.CreatingTaskZset, &z)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add task to creating task list: ", taskId, cmd.Err().Error())
		return err
	}

	glog.Info("succeeded to add task to creating task list: ", taskId)
	return nil
}

// 将任务添加到任务列表中
// 将 $taskid 推入 redis_task_zset。表示任务已经存在。
func AddTaskToExistingTaskList(taskId taskmodel.TaskIdType, timeout time.Duration) error {

	var z = redis.Z{
		Score:  float64(time.Now().Add(timeout).Unix()),
		Member: taskId,
	}

	cmd := redistool.DefaultRedis().ZAdd(context.Background(), config.TaskZset, &z)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add task to existing task list: ", taskId, err.Error())
		return err
	}

	glog.Info("succeeded to add task to existing task list: ", taskId)
	return nil
}

// 将任务推送到已完成队列, 等待任务管理逻辑进行处理
func PushTaskToCompletedList(taskId taskmodel.TaskIdType) error {

	z := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: taskId,
	}

	cmd := redistool.DefaultRedis().ZAdd(context.Background(), config.CompletedTaskList, &z)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to add task to completed list: ", taskId, err.Error())
		return err
	}

	glog.Info("succeeded to push task to completed list: ", taskId)
	return nil
}
