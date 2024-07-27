package main

import (
	"fmt"
	"sync"
)

type Storer interface {
	Push([]byte) (int, error)
	Pull(int) ([]byte, error)
	Len() int
}

type StoreProducerFunc func() Storer

type MemoryStore struct {
	mu   sync.RWMutex
	data [][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make([][]byte, 0),
	}
}

func (s *MemoryStore) Push(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

func (s *MemoryStore) Pull(offset int) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be less then zero")
	}
	if len(s.data) < offset {
		return nil, fmt.Errorf("offest (%d) too high", offset)
	}
	return s.data[offset], nil
}

func (s *MemoryStore) Len() int {
	return len(s.data)
}
