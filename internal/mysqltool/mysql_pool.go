package mysqltool

import (
	// 引入Go SQL驱动
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"

	"pterergate-dtf/dtf/extconfig"
	"pterergate-dtf/internal/config"
)

// 默认的MySQL数据库对象
var gs_MySQLDB *sqlx.DB

// 获取默认的MySQL实例
func DefaultMySQL() *sqlx.DB {
	return gs_MySQLDB
}

// 连接默认的MySQL数据库
func ConnectToDefaultMySQL() {
	InitMySQLClient(&config.DefaultMySQL, "", &gs_MySQLDB)
}

// 初始化MySQL连接
func InitMySQLClient(sqlConf *extconfig.MySQLAddress, setting string, targetDB **sqlx.DB) {

	if len(setting) > 0 {
		setting = fmt.Sprintf("&%s", setting)
	}

	addr := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4%s",
		sqlConf.Username,
		sqlConf.Password,
		sqlConf.Protocol,
		sqlConf.Address,
		sqlConf.DB,
		setting,
	)

	// 连接数据库，并发送ping进行验证，保证连接成功
	db, err := sqlx.Connect(sqlConf.Type, addr)
	if err != nil {
		glog.Warning("failed to connect to MySQL: ", err, sqlConf)
		panic(err)
	}

	// 设置连接池配置
	db.SetConnMaxIdleTime(config.MySQL_ConnMaxIdleTime * time.Second)
	db.SetConnMaxLifetime(config.MySQL_ConnMaxLifeTime * time.Second)
	db.SetMaxOpenConns(config.MySQL_MaxOpenConns)
	db.SetMaxIdleConns(config.MySQL_MaxIdleConns)

	// 初始只有一个连接，这里创建多个连接
	for i := 0; i < config.MySQL_InitialOpenConnections; i++ {
		ctx := context.Background()
		conn, err := db.Conn(ctx)
		if err != nil {
			glog.Warning("failed to open new connection: ", err)
		}

		defer conn.Close()
	}

	*targetDB = db
}
