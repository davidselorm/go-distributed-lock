package distlock

import (
	"testing"
	"time"
)

func TestLockAcquireAndRelease(t *testing.T) {
	lock := NewLock("database-cluster", 1*time.Second)

	tok1, err := lock.TryAcquire("node-A")
	if err != nil || tok1 == 0 {
		t.Fatalf("node-A should acquire lock, got token: %d, err: %v", tok1, err)
	}

	_, err = lock.TryAcquire("node-B")
	if err != ErrLockHeld {
		t.Fatalf("node-B should fail while node-A holds lock, got %v", err)
	}

	if !lock.Release("node-A") {
		t.Fatalf("node-A should release lock")
	}

	tok2, err := lock.TryAcquire("node-B")
	if err != nil || tok2 <= tok1 {
		t.Fatalf("node-B should acquire with monotonically increasing token, got: %d", tok2)
	}
}
