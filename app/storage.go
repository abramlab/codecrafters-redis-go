package main

import "sync"

type storage struct {
	data map[string]string
	mu   sync.RWMutex
}

func newStorage() *storage {
	return &storage{
		data: make(map[string]string),
	}
}

func (s *storage) get(key string) (string, bool) {
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()
	return val, ok
}

func (s *storage) set(key, val string) {
	s.mu.Lock()
	s.data[key] = val
	s.mu.Unlock()
}
