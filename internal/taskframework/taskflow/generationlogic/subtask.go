package generationlogic

import (
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/config"
	"github.com/danenmao/pterergate-dtf/internal/idtool"
)

const (
	SubtaskDefaultTimeout = 300
	SubtaskMaxTimeout     = 1800
)

// 创建一个子任务
func CreateSubtask(
	taskId taskmodel.TaskIdType,
	taskType uint32,
	impl taskmodel.ITaskGenerator,
	subtaskData *taskmodel.SubtaskBody,
	finished *bool,
) error {

	// invoke the plugin generator to get a subtask data
	subtaskStartTime := time.Now()
	err := impl.GetSubtask(taskId, subtaskData, finished)

	costTime := time.Since(subtaskStartTime)
	if costTime > time.Second*20 {
		glog.Warningf("GetSubtask %d costs too much time: %ds", taskId, costTime/time.Second)
	}

	if err == errordef.ErrNotFound {
		return err
	}

	if err != nil {
		glog.Warning("failed to invoke subtask fn: ", taskId, ",", err.Error())
		return err
	}

	// get a subtask id
	id, err := idtool.GetId(config.SubtaskIdKey)
	if err != nil {
		glog.Warning("failed to get a subtask id: ", err)
		return err
	}

	subtaskId := taskmodel.SubtaskIdType(id)

	// set subtask data
	subtaskData.SubtaskId = subtaskId
	subtaskData.TaskId = taskId
	subtaskData.TaskType = taskType
	subtaskData.CreatedAt = time.Now()

	// check if the subtask timeout is valid
	if subtaskData.Timeout <= 0 {
		subtaskData.Timeout = SubtaskDefaultTimeout
	} else if subtaskData.Timeout > SubtaskMaxTimeout {
		subtaskData.Timeout = SubtaskMaxTimeout
	}

	glog.Info("succeeded to create a subtask: ", subtaskId, taskId)
	return nil
}
