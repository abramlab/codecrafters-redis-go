package main

import (
	"sync"
	"time"
)

var nowTime = time.Now

type value struct {
	v   string
	exp time.Time
}

type storage struct {
	data map[string]*value
	mu   sync.RWMutex
}

func newStorage() *storage {
	return &storage{
		data: make(map[string]*value),
	}
}

func (s *storage) get(key string) (string, bool) {
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return "", false
	}

	if val.exp.IsZero() || nowTime().Before(val.exp) {
		return val.v, ok
	}

	// remove expired key and return empty string
	s.remove(key)
	return "", false
}

func (s *storage) set(key, val string) {
	s.mu.Lock()
	s.data[key] = &value{v: val}
	s.mu.Unlock()
}

func (s *storage) remove(key string) {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
}

func (s *storage) setWithExpiration(key, val string, exp time.Duration) {
	s.mu.Lock()
	s.data[key] = &value{v: val, exp: nowTime().Add(exp)}
	s.mu.Unlock()
}
