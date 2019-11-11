package kv

import "sync"

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

type MapStore struct {
	Map map[string]string
	mu  sync.RWMutex
}

func (m MapStore) Set(key string, value string) error {
	m.mu.Lock()
	m.Map[key] = value
	m.mu.Unlock()
	return nil
}

func (m MapStore) Get(key string) (string, error) {
	m.mu.RLock()
	value := m.Map[key]
	m.mu.RUnlock()
	return value, nil
}
