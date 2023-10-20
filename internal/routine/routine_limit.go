package routine

import (
	"sync"
)

// 例程限制结构
type RoutineCountLimiter struct {
	UpperLimit uint32     // 上限值
	counter    uint32     // 当前实例中正在执行的例程数
	lock       sync.Mutex // 锁
}

// 返回当前的例程数
func (limit *RoutineCountLimiter) GetCount() uint32 {
	return limit.counter
}

// 检查当前实例的例程数是否超过上限
func (limit *RoutineCountLimiter) IsFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()
	return limit.counter >= limit.UpperLimit
}

// 如果当前实例的例程数未超过上限，增加计数
// 返回是否成功增加计数
func (limit *RoutineCountLimiter) IncrIfNotFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()

	if limit.counter >= limit.UpperLimit {
		return false
	}

	// 增加计数
	limit.counter += 1
	return true
}

// 增加正在执行的例程数
func (limit *RoutineCountLimiter) Incr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.counter += 1
}

// 减少正在执行的例程数
func (limit *RoutineCountLimiter) Decr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.counter -= 1
}
