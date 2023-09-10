package tasktool

import (
	"fmt"

	"github.com/golang/glog"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/dbdef"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/taskframework/taskflow/flowdef"
)

// 添加任务记录
func AddTaskRecord(task *dbdef.TaskRecord) error {

	result, err := mysqltool.DefaultMySQL().NamedExec(
		dbdef.SQL_TaskTable_InsertTask,
		task,
	)

	if err != nil {
		glog.Warning("failed to add task record: ", err.Error())
		return err
	}

	lines, _ := result.RowsAffected()
	glog.Info("added a task record: ", task.Id, lines)

	return nil
}

// 获取task info key的名称
func GetTaskInfoKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", config.TaskInfoKeyPrefix, taskId)
}

// 获取保存任务调度数据的Key
func GetTaskCreateParamKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", flowdef.TaskCreateParamPrefix, taskId)
}

func GetTaskLockKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", config.TaskInfoLockPrefix, taskId)
}

func GetTaskGenerationProgressKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", config.TaskGenerationKeyPrefix, taskId)
}

func GetTaskSubtaskListKey(taskId taskmodel.TaskIdType) string {
	return fmt.Sprintf("%s%d", config.TaskToSubtaskSetPrefix, taskId)
}
