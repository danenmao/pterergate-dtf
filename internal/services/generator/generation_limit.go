package generator

import (
	"pterergate-dtf/internal/misc"
	"pterergate-dtf/internal/routine"
)

const (
	GenerationRoutineCountDefaultLimit = 200
	GeneratingRoutineLimitEnvName      = "GENERATION_ROUTINE_LIMIT"
)

var (
	// 当前实例中并行执行生成的任务数上限
	gs_GenerationRoutineCountLimit = routine.RoutineCountLimit{}
)

// 初始化
func init() {
	gs_GenerationRoutineCountLimit.RoutineCountLimit = uint32(misc.GetIntFromEnv(GeneratingRoutineLimitEnvName,
		GenerationRoutineCountDefaultLimit))
}

// 检查当前实例生成的任务数是否超过上限
func CheckIfExceedLimit() bool {
	return gs_GenerationRoutineCountLimit.CheckIfExceedLimit()
}

// 如果当前实例生成的任务数未超过上限，增加计数
// 返回是否成功增加计数
func IncrIfNotExceedLimit() bool {
	return gs_GenerationRoutineCountLimit.IncrIfNotExceedLimit()
}

// 增加正在生成的例程数
func IncrGeneratingRoutineCount() {
	gs_GenerationRoutineCountLimit.IncrRoutineCount()
}

// 减少正在生成的例程数
func DecrGeneratingRoutineCount() {
	gs_GenerationRoutineCountLimit.DecrRoutineCount()
}
