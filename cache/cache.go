package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiryTime time.Time
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]cacheItem),
		mu:    sync.RWMutex{},
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheItem{
		value:      value,
		expiryTime: time.Now().UTC().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cacheItem, ok := c.cache[key]
	if !ok || cacheItem.expiryTime.Before(time.Now().UTC()) {
		return nil, false
	}

	return cacheItem.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
