package schedulerflow

import (
	"github.com/golang/glog"

	"pterergate-dtf/dtf/errordef"
	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/taskframework/taskflow/subtaskqueue"
	"pterergate-dtf/internal/tasktool"
)

// get a subtask from the subtask queue
func GetSubtask(
	taskId taskmodel.TaskIdType,
	subtaskData *taskmodel.SubtaskData,
	finished *bool,
) error {

	// get a subtask from the subtask queue
	queue := subtaskqueue.SubtaskQueue{TaskId: taskId}
	err := queue.PopSubtask(subtaskData)
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
