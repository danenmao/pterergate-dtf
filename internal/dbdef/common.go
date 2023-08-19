package dbdef

import "time"

const (
	Query_DefaultOffset = 0
	Query_DefaultLimit  = 10
	Query_MaxLimit      = 100
)

const (
	GoTimeFormatStr = "2006-01-02 15:04:05"
	DBNullTimeStr   = "1900-01-01 00:00:00"
	DBListSeperator = ";"

	SQL_SelectFoundRows = "SELECT FOUND_ROWS() as total"
)

// 获取默认空时间
func GetDBNullTime() time.Time {
	loc, _ := time.LoadLocation("Local")
	DBNullTime, _ := time.ParseInLocation(GoTimeFormatStr, DBNullTimeStr, loc)
	return DBNullTime
}

// 获取指定时间
func GetTimeInLocal(timestr string) time.Time {
	loc, _ := time.LoadLocation("Local")
	goalTime, _ := time.ParseInLocation(GoTimeFormatStr, timestr, loc)
	return goalTime
}
