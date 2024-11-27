package cache

import (
	"sync"
	"time"
)

type Cache struct {
	lru *LRU
	mu  sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewCache(capacity int) *Cache {
	c := &Cache{
		lru: NewLRU(capacity),
		mu:  sync.RWMutex{},
	}

	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.lru.Set(key, cacheItem{
		value:      value,
		expiration: expiration,
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()

	item, ok := c.lru.Get(key)
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}

	cacheItem := item.(cacheItem)
	if !cacheItem.expiration.IsZero() && cacheItem.expiration.Before(time.Now()) {
		c.mu.RUnlock()
		c.Delete(key)
		return nil, false
	}

	c.mu.RUnlock()
	return cacheItem.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Delete(key)
}
