package config

import (
	"time"

	"pterergate-dtf/internal/basedef"
)

// manager 使用的设置
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
