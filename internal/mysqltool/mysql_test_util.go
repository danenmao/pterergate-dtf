package mysqltool

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var s_actualDB *sql.DB
var s_mockDB *sql.DB
var DBMock sqlmock.Sqlmock

// 配置测试环境
func Setup() {

	if DefaultMySQL() != nil {
		s_actualDB = DefaultMySQL().DB
	}

	var err error
	s_mockDB, DBMock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic("failed to mock sql")
	}

	gs_MySQLDB = sqlx.NewDb(s_mockDB, "mysql")
}

// 清理测试环境
func Teardown() {

	if s_actualDB != nil {
		DefaultMySQL().DB = s_actualDB
		s_actualDB = nil
	}

	if gs_MySQLDB != nil {
		gs_MySQLDB.Close()
		gs_MySQLDB = nil
	}

	if s_mockDB != nil {
		s_mockDB.Close()
		s_mockDB = nil
	}
}
