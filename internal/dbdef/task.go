package dbdef

import "fmt"

// 任务结构
type TaskRecord struct {
	Id            uint64 `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	Description   string `db:"description" json:"description"`
	UID           uint64 `db:"uid" json:"uid"`
	Creator       string `db:"creator" json:"creator"`
	StartTime     string `db:"start_time" json:"start_time"`
	FinishTime    string `db:"finish_time" json:"finish_time"`
	NextCheckTime string `db:"next_check_time" json:"next_check_time"`
	TaskType      uint32 `db:"task_type" json:"task_type"`
	TimeCost      uint32 `db:"time_cost" json:"time_cost"`
	TaskStatus    uint8  `db:"task_status" json:"task_status"`
}

// 任务表的定义
const (
	TaskTableName           = "tbl_tcss_compliance_task"
	TaskTable_Id            = "id"
	TaskTable_Name          = "name"
	TaskTable_Description   = "description"
	TaskTable_UID           = "uid"
	TaskTable_Creator       = "creator"
	TaskTable_StartTime     = "start_time"
	TaskTable_FinishTime    = "finish_time"
	TaskTable_NextCheckTime = "next_check_time"
	TaskTable_TaskType      = "task_type"
	TaskTable_TimeCost      = "time_cost"
	TaskTable_TaskStatus    = "task_status"
)

// 创建任务表的语句
var SQL_CreateTaskTable string = fmt.Sprintf(
	"CREATE TABLE IF NOT EXISTS `%s` ("+
		"`id` bigint UNSIGNED NOT NULL,"+
		"`name` varchar(100) NOT NULL COMMENT '任务名称',"+
		"`description` varchar(255) NOT NULL DEFAULT '' COMMENT '任务的描述',"+
		"`uid` bigint UNSIGNED NOT NULL COMMENT '所有者id',"+
		"`creator` varchar(100) NOT NULL COMMENT '创建者的账号',"+
		"`start_time` datetime NOT NULL COMMENT '创建任务的时间',"+
		"`finish_time` datetime NOT NULL COMMENT '任务结束的时间',"+
		"`next_check_time` datetime NOT NULL COMMENT '下次检查任务状态的时间',"+
		"`task_type` int UNSIGNED NOT NULL COMMENT '任务类型',"+
		"`time_cost` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务耗时,单位为秒',"+
		"`task_status` tinyint UNSIGNED NOT NULL COMMENT '任务的状态',"+

		"PRIMARY KEY (`id`),"+
		"KEY `key_uid` (`uid`,`asset_type`,`risk_level`),"+
		"KEY `key_start_time` (`start_time`, `finish_time`)"+
		")"+
		"ENGINE = InnoDB "+
		"AUTO_INCREMENT = 1 "+
		"DEFAULT CHARSET = utf8mb4 "+
		"COMMENT='任务表'",

	TaskTableName,
)

// 添加任务记录
var SQL_TaskTable_InsertTask string = fmt.Sprintf(
	"INSERT INTO `%s` (`%s`,`%s`,`%s`,`%s`,`%s`,`%s`,`%s`,`%s`,`%s`,`%s`,"+
		") VALUES (:%s,:%s,:%s,:%s,:%s,:%s,:%s,:%s,:%s)",

	TaskTableName,

	TaskTable_Id,
	TaskTable_Name,
	TaskTable_Description,
	TaskTable_UID,
	TaskTable_Creator,
	TaskTable_StartTime,
	TaskTable_FinishTime,
	TaskTable_NextCheckTime,
	TaskTable_TaskType,
	TaskTable_TimeCost,

	TaskTable_Id,
	TaskTable_Name,
	TaskTable_Description,
	TaskTable_UID,
	TaskTable_Creator,
	TaskTable_StartTime,
	TaskTable_FinishTime,
	TaskTable_NextCheckTime,
	TaskTable_TaskType,
)

// 任务完成中更新任务记录
var SQL_TaskTable_CompleteTask string = fmt.Sprintf(
	"UPDATE `%s` SET `%s`=:%s,`%s`=:%s,`%s`=:%s where `%s`=:%s",
	TaskTableName,
	TaskTable_FinishTime,
	TaskTable_FinishTime,
	TaskTable_TimeCost,
	TaskTable_TimeCost,
	TaskTable_TaskStatus,
	TaskTable_TaskStatus,
	TaskTable_Id,
	TaskTable_Id,
)

// 任务执行过程中，更新任务的检查时间
var SQL_TaskTable_UpdateNextCheckTime string = fmt.Sprintf(
	"UPDATE `%s` SET `%s`=? where `%s`=?",
	TaskTableName,
	TaskTable_NextCheckTime,
	TaskTable_Id,
)

// 获取任务表中创建异常的任务列表
var SQL_TaskTable_QueryExceptionalCreationTask string = fmt.Sprintf(
	"select `%s` from `%s` where `%s` < ? limit ?, ? ",
	TaskTable_Id,
	TaskTableName,
	TaskTable_NextCheckTime,
)
