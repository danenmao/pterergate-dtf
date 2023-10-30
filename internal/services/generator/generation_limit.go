package generator

import (
	"github.com/danenmao/pterergate-dtf/internal/misc"
	"github.com/danenmao/pterergate-dtf/internal/routine"
)

const (
	GenerationCountDefaultLimit = 200
	GenerationCountLimitEnv     = "GENERATION_ROUTINE_LIMIT"
)

var (
	// 当前实例中并行执行生成的任务数上限
	gs_GenerationLimiter = routine.CountLimiter{}
)

// 初始化
func init() {
	gs_GenerationLimiter.UpperLimit = uint32(misc.GetIntFromEnv(GenerationCountLimitEnv,
		GenerationCountDefaultLimit))
}

// 检查当前实例生成的任务数是否超过上限
func IsFull() bool {
	return gs_GenerationLimiter.IsFull()
}

// 如果当前实例生成的任务数未超过上限，增加计数
// 返回是否成功增加计数
func IncrIfNotFull() bool {
	return gs_GenerationLimiter.IncrIfNotFull()
}

// 增加正在生成的例程数
func Incr() {
	gs_GenerationLimiter.Incr()
}

// 减少正在生成的例程数
func Decr() {
	gs_GenerationLimiter.Decr()
}
