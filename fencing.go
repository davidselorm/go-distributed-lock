package lock

import "sync/atomic"

var globalFence int64

func (l *DistLock) AcquireWithFence() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FencingToken = atomic.AddInt64(&globalFence, 1)
	return l.FencingToken
}
