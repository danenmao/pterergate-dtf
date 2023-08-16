package config

// MySQL数据库的连接配置
var MySQLConf = map[string]string{
	"name":     "",
	"type":     "",
	"protocol": "",
	"address":  "",
	"db":       "",
	"username": "",
	"auth":     "",
}

// MySQL连接配置
const (
	MySQL_InitialOpenConnections = 10
	MySQL_MaxOpenConns           = 100
	MySQL_MaxIdleConns           = 10
	MySQL_ConnMaxIdleTime        = 900
	MySQL_ConnMaxLifeTime        = 2 * 3600
)
