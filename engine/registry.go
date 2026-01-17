package engine

import (
	"sync"
	"io"
)

type Registry struct {
	data    map[string][]byte
	readers map[string]io.Reader
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		data:    make(map[string][]byte),
		readers: make(map[string]io.Reader),
	}
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

func (r *Registry) SetReader(key string, reader io.Reader) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readers[key] = reader
}

func (r *Registry) GetReader(key string) io.Reader {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.readers[key]
}
