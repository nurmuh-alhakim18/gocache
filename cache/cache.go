package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheItem
	mu       sync.RWMutex
	stopChan chan struct{}
}

type cacheItem struct {
	value      interface{}
	expiryTime time.Time
}

func NewCache() *Cache {
	c := &Cache{
		cache:    make(map[string]cacheItem),
		mu:       sync.RWMutex{},
		stopChan: make(chan struct{}),
	}

	go c.cleanUp()
	return c
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
	cacheItem, ok := c.cache[key]
	c.mu.RUnlock()

	if !ok || cacheItem.expiryTime.Before(time.Now().UTC()) {
		c.Delete(key)
		return nil, false
	}

	return cacheItem.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *Cache) cleanUp() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.cache {
				if item.expiryTime.Before(time.Now().UTC()) {
					delete(c.cache, key)
				}
			}

			c.mu.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.stopChan)
}
