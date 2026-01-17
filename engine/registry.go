package engine

import "sync"

type Registry struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{data: make(map[string][]byte)}
}

func (r *Registry) Set(key string, value []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

func (r *Registry) Get(key string) []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[key]
}
