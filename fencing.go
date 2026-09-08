package distlock

import (
	"sync/atomic"
)

type FencedStorage struct {
	lastToken atomic.Uint64
	store     map[string][]byte
}

func NewFencedStorage() *FencedStorage {
	return &FencedStorage{
		store: make(map[string][]byte),
	}
}

func (fs *FencedStorage) Put(fenceToken uint64, key string, val []byte) error {
	prev := fs.lastToken.Load()
	if fenceToken < prev {
		return ErrInvalidFence
	}
	fs.lastToken.Store(fenceToken)
	fs.store[key] = val
	return nil
}
