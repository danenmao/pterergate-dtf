package mysqltool

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var s_actualDB *sqlx.DB
var s_mockDB *sqlx.DB
var DBMock sqlmock.Sqlmock

// 配置测试环境
func Setup() {
	// save actual sqlx.DB
	if gs_MySQLDB != nil {
		s_actualDB = gs_MySQLDB
	}

	// generate a mock DB
	var err error
	var mockSQL *sql.DB
	mockSQL, DBMock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic("failed to mock sql")
	}

	s_mockDB = sqlx.NewDb(mockSQL, "mysql")

	// overwrite default sqlx.DB
	gs_MySQLDB = s_mockDB
}

// 清理测试环境
func Teardown() {
	// restore default sqlx.DB
	if s_actualDB != nil {
		gs_MySQLDB = s_actualDB
		s_actualDB = nil
	}

	// release mock sqlx.DB
	if s_mockDB != nil {
		s_mockDB.Close()
		s_mockDB = nil
	}
}
