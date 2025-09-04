package cache

import "sync"

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func New() *MemoryCache {
	return &MemoryCache{data: make(map[string]interface{})}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *MemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
