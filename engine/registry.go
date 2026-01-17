package engine

type Registry struct {
	data map[string][]byte
}

func NewRegistry() *Registry {
	return &Registry{data: make(map[string][]byte)}
}

func (r *Registry) Set(key string, value []byte) {
	r.data[key] = value
}

func (r *Registry) Get(key string) []byte {
	return r.data[key]
}
