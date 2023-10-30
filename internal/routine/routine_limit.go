package routine

import (
	"sync"
)

// count limiter
type CountLimiter struct {
	UpperLimit uint32     // the upper limit
	counter    uint32     // the current count
	lock       sync.Mutex // the lock
}

// return the current count
func (limit *CountLimiter) Count() uint32 {
	return limit.counter
}

// check if the count exceeds the upper limit
func (limit *CountLimiter) IsFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()
	return limit.counter >= limit.UpperLimit
}

// if the count don't exceed the upper limit, incr the count
// return if incr the count
func (limit *CountLimiter) IncrIfNotFull() bool {

	limit.lock.Lock()
	defer limit.lock.Unlock()

	if limit.counter >= limit.UpperLimit {
		return false
	}

	limit.counter += 1
	return true
}

// incr the count
func (limit *CountLimiter) Incr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.counter += 1
}

// decr the count
func (limit *CountLimiter) Decr() {
	limit.lock.Lock()
	defer limit.lock.Unlock()
	limit.counter -= 1
}
