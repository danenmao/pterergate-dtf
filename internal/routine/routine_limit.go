package routine

import (
	"sync"
)

// 例程限制结构
type RoutineCountLimit struct {
	RoutineCountLimit   uint32     // 上限值
	currentRoutineCount uint32     // 当前实例中正在执行的例程数
	routineCountLock    sync.Mutex // 锁
}

// 返回当前的例程数
func (limit *RoutineCountLimit) GetRoutineCount() uint32 {
	return limit.currentRoutineCount
}

// 检查当前实例的例程数是否超过上限
func (limit *RoutineCountLimit) CheckIfExceedLimit() bool {

	limit.routineCountLock.Lock()
	defer limit.routineCountLock.Unlock()
	return limit.currentRoutineCount >= limit.RoutineCountLimit
}

// 如果当前实例的例程数未超过上限，增加计数
// 返回是否成功增加计数
func (limit *RoutineCountLimit) IncrIfNotExceedLimit() bool {

	limit.routineCountLock.Lock()
	defer limit.routineCountLock.Unlock()

	if limit.currentRoutineCount >= limit.RoutineCountLimit {
		return false
	}

	// 增加计数
	limit.currentRoutineCount += 1
	return true
}

// 增加正在执行的例程数
func (limit *RoutineCountLimit) IncrRoutineCount() {
	limit.routineCountLock.Lock()
	defer limit.routineCountLock.Unlock()
	limit.currentRoutineCount += 1
}

// 减少正在执行的例程数
func (limit *RoutineCountLimit) DecrRoutineCount() {
	limit.routineCountLock.Lock()
	defer limit.routineCountLock.Unlock()
	limit.currentRoutineCount -= 1
}
