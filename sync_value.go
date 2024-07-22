package main

import (
	"sync"
)

type BoolValue struct {
	value bool
	mu    sync.RWMutex
}

func NewBoolValue(value bool) *BoolValue {
	return &BoolValue{value: value, mu: sync.RWMutex{}}
}

func (b *BoolValue) Set(value bool) {
	b.mu.Lock()
	b.value = value
	b.mu.Unlock()
}

func (b *BoolValue) Get() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.value
}
