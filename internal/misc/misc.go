package misc

import (
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"pterergate-dtf/internal/config"
)

// 获取传入的函数的名称，格式为 package/package.name
func GetFunctionName(i interface{}, seps ...rune) string {

	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}

	return ""
}

// 从环境变量得到整数值
func GetIntFromEnv(env string, defaultVal int) int {

	val := os.Getenv(env)
	if len(val) <= 0 {
		return defaultVal
	}

	retVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return retVal
}

// 仅在测试模式下输出内容
func DumpDataInTest(prompt string, in interface{}) {

	// 仅在dev, test上执行
	if config.WorkEnv != config.ENV_DEV && config.WorkEnv != config.ENV_TEST {
		return
	}

	data, err := json.Marshal(in)
	if err != nil {
		glog.Warning("failed to marshal object: ", err)
		return
	}

	glog.Info(prompt, ": ", string(data))
}
