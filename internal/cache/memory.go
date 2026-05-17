package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]entry),
		ttl:   ttl,
	}

	go c.cleanupLoop()

	return c
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(e.expiresAt) {
		return nil, false
	}

	return e.value, true
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.evict()
	}
}

func (c *Cache) evict() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, key)
		}
	}
}
