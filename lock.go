package lock

import (
	"sync"
	"time"
)

type DistLock struct {
	mu sync.Mutex
	Key string
	Owner string
	TTL time.Duration
	FencingToken int64
}

func NewLock(key, owner string, ttl time.Duration) *DistLock {
	return &DistLock{Key: key, Owner: owner, TTL: ttl, FencingToken: 1}
}
