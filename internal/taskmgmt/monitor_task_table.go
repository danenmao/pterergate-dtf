package taskmgmt

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"

	"pterergate-dtf/dtf/taskmodel"
	"pterergate-dtf/internal/basedef"
	"pterergate-dtf/internal/config"
	"pterergate-dtf/internal/dbdef"
	"pterergate-dtf/internal/mysqltool"
	"pterergate-dtf/internal/redistool"
	"pterergate-dtf/internal/taskframework/taskflow/flowdef"
	"pterergate-dtf/internal/tasktool"
)

// 监视任务表中的任务是否正常
func MonitorTaskTableRoutine() {
	// 从 tbl_compliance_task 取过了检查时间的任务的记录
	var taskList []taskmodel.TaskIdType
	getExceptionalTasks(&taskList)

	// 处理异常的任务
	for _, taskId := range taskList {
		repairExceptionalTask(taskId)
	}
}

// 获取异常创建的任务的列表
func getExceptionalTasks(taskList *[]taskmodel.TaskIdType) error {

	queryFn := func(offset int, limit int) (*sqlx.Rows, error) {
		return mysqltool.DefaultMySQL().Queryx(
			dbdef.SQL_TaskTable_QueryExceptionalCreationTask,
			time.Now().Format(basedef.GoTimeFormatStr),
			offset, limit,
		)
	}

	readFn := func(rows *sqlx.Rows) error {
		var id uint64 = 0
		err := rows.Scan(&id)
		if err != nil {
			glog.Warning("failed to get exceptional task id: ", err)
			return nil
		}

		*taskList = append(*taskList, taskmodel.TaskIdType(id))
		return nil
	}

	err := mysqltool.ReadFromDBByPageCustom(queryFn, readFn, 10)
	if err != nil {
		glog.Warning("failed to get exceptional tasks: ", err)
		return err
	}

	glog.Info("succeeded to get exceptional tasks: ", taskList)
	return nil
}

// 修复创建过程异常的任务
func repairExceptionalTask(taskId taskmodel.TaskIdType) {

	cmd := redistool.DefaultRedis().ZScore(context.Background(), config.CreatingTaskZset,
		strconv.FormatUint(uint64(taskId), 10))

	// zscore，取一个不存在的key, go得到的是err(redis: nil)，冒号后面有个空格
	// zscore，取一个不存在的member时, go得到的是err(redis: nil)，冒号后面有个空格
	if cmd.Err() == nil {
		glog.Info("need to do nothing for task: ", taskId)
		return
	}

	glog.Info("try to repair task creation: ", taskId)

	// 重新获取任务结构
	var taskParam = taskmodel.TaskParam{}
	err := RefillTaskParam(taskId, &taskParam)
	if err != nil {
		glog.Warning("failed to refill task record for task: ", taskId, err)
		return
	}

	// 修复任务
	go TaskCreationRoutine(taskId, 0, &taskParam)
}

// 重新填写任务结构
func RefillTaskParam(
	taskId taskmodel.TaskIdType,
	taskParam *taskmodel.TaskParam,
) error {

	if taskId == 0 {
		glog.Warning("invalid task id")
		return errors.New("invalid task id")
	}

	createParam := flowdef.TaskCreateParam{}
	err := tasktool.GetTaskCreateParam(taskId, &createParam)
	if err != nil {
		glog.Warning("failed to get the create param of task: ", taskId, err)
		return err
	}

	var typeParam string
	err = tasktool.GetTaskRawTypeParam(taskId, &typeParam)
	if err != nil {
		glog.Warning("failed to get the type param of task: ", taskId, err)
		return err
	}

	glog.Info("succeeded to get init task record of task: ", taskId)
	return nil
}
