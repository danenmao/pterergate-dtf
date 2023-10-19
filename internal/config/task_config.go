package config

import (
	"time"

	"github.com/danenmao/pterergate-dtf/internal/basedef"
)

// manager settings
var (
	//
	// go_task_creation的设置
	//
	EnvTaskNextCheckTimeMax  string = time.Now().AddDate(3000, 1, 1).Format(basedef.GoTimeFormatStr)
	EnvTaskCreationNextCheck int    = 120
	EnvTaskCreatingTimeout   int    = 100

	//
	// go_monitor_task_tbl的设置
	//
	EnvMonitorTaskTblCountLimit   uint = 2
	EnvMonitorTaskTblInterval     int  = 120
	EnvMonitorTaskTbleRecordLimit uint = 1000

	//
	// go_monitor_task_creation的设置
	//
	EnvMonitorTaskCreationCountLimit uint = 2
	EnvMonitorTaskCreationInterval   int  = 60

	//
	// go_monitor_task_timeout的设置
	//
	EnvMonitorTaskTimeoutCountLimit uint = 2
	EnvMonitorTaskTimeoutInterval   int  = 30
	EnvTaskTimeout                  int  = 1800

	//
	// go_monitor_completed_task的设置
	//
	EnvMonitorTaskCompletedCountLimit uint = 2
	EnvMonitorTaskCompletedInterval   int  = 5
)

// generator settings
var (

	// go_start_generate_task
	EnvGenerateTaskConcurrencyLimit uint = 5
	EnvGenerateTaskCheckInterval    int  = 1

	// go_monitor_task_generation
	EnvMonitorTaskGenerationConcurrencyLimit uint = 2
	EnvMonitorTaskGenerationInterval         int  = 30
)

// scheduler settings
var (

	// go_schedule_task
	EnvScheduleTaskConcurrencyLimit uint = 10
	EnvScheduleTaskInterval         int  = 200

	// go_retry_push_subtask
	EnvRetryPushSubtaskConcurrencyLimit uint = 1
	EnvRetryPushSubtaskInterval         int  = 2

	// go_monitor_subtask_timeout
	EnvMonitorSubtaskTimeoutConcurrencyLimit uint = 5
	EnvMonitorSubtaskTimeoutInterval              = 2

	// go_monitor_subtask_complete
	EnvMonitorSubtaskCompleteConcurrencyLimit uint = 10
	EnvMonitorSubtaskCompleteInterval         int  = 200

	// go_monitor_task_complete
	EnvMonitorTaskCompleteConcurrencyLimit uint = 2
	EnvMonitorTaskCompleteInterval         int  = 1
)

// collector settings
var (
	//
	// go_complete_subtask
	//
	EnvCompleteSubtaskConcurrencyLimit uint = 2
	EnvCompleteSubtaskInterval         int  = 100
)
