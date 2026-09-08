package distlock

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrLockHeld = errors.New("lock already held")
	ErrLockTimeout = errors.New("timed out waiting for lock")
	ErrInvalidFence = errors.New("fencing token validation failed")
)

type Lock struct {
	mu           sync.Mutex
	resource     string
	owner        string
	tokenCounter atomic.Uint64
	currentToken uint64
	leaseExpiry  time.Time
	leaseTTL     time.Duration
}

func NewLock(resource string, ttl time.Duration) *Lock {
	return &Lock{
		resource: resource,
		leaseTTL: ttl,
	}
}

func (l *Lock) TryAcquire(owner string) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if l.owner != "" && now.Before(l.leaseExpiry) && l.owner != owner {
		return 0, ErrLockHeld
	}

	l.owner = owner
	l.leaseExpiry = now.Add(l.leaseTTL)
	l.currentToken = l.tokenCounter.Add(1)
	return l.currentToken, nil
}

func (l *Lock) AcquireWithContext(ctx context.Context, owner string, retryInterval time.Duration) (uint64, error) {
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		token, err := l.TryAcquire(owner)
		if err == nil {
			return token, nil
		}

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (l *Lock) Release(owner string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.owner == owner {
		l.owner = ""
		l.leaseExpiry = time.Time{}
		return true
	}
	return false
}

func (l *Lock) ValidateFencingToken(token uint64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return token == l.currentToken && time.Now().Before(l.leaseExpiry)
}
