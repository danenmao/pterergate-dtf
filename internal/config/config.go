package config

var (
	// 工作环境定义
	WorkEnv uint16

	// 配置文件所在目录
	ConfigDir string
)

// 工作环境值定义
const (
	ENV_DEV       = 1
	ENV_TEST      = 2
	ENV_PREONLINE = 3
	ENV_ONLINE    = 4

	ENV_DEV_STR       = "DEV"
	ENV_TEST_STR      = "TEST"
	ENV_PREONLINE_STR = "PREONLINE"
	ENV_ONLINE_STR    = "ONLINE"
)
