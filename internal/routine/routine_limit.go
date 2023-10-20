package routine

import (
	"sync"
)

// 例程限制结构
type RoutineCountLimit struct {
	CountLimit   uint32     // 上限值
	currentCount uint32     // 当前实例中正在执行的例程数
	lock         sync.Mutex // 锁
}

// 返回当前的例程数
func (limit *RoutineCountLimit) GetCount() uint32 {
	return limit.currentCount
}

// 检查当前实例的例程数是否超过上限
func (limit *RoutineCountLimit) IsFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()
	return limit.currentCount >= limit.CountLimit
}

// 如果当前实例的例程数未超过上限，增加计数
// 返回是否成功增加计数
func (limit *RoutineCountLimit) IncrIfNotFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()

	if limit.currentCount >= limit.CountLimit {
		return false
	}

	// 增加计数
	limit.currentCount += 1
	return true
}

// 增加正在执行的例程数
func (limit *RoutineCountLimit) Incr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.currentCount += 1
}

// 减少正在执行的例程数
func (limit *RoutineCountLimit) Decr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.currentCount -= 1
}
