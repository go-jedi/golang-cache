package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type cacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

type Cache struct {
	m          sync.RWMutex
	data       map[string]cacheItem
	defaultTTL time.Duration
}

// NewCache initialize cache.
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}
}

// Set data by key in cache.
func (c *Cache) Set(key string, data interface{}, ttl ...time.Duration) {
	c.m.Lock()
	defer c.m.Unlock()

	duration := c.defaultTTL
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	c.data[key] = cacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

// Get data by key from cache.
func (c *Cache) Get(key string, v interface{}) bool {
	c.m.RLock()
	defer c.m.RUnlock()

	item, ok := c.data[key]
	if !ok || time.Now().After(item.ExpiresAt) {
		return false
	}

	dataBytes, err := json.Marshal(item.Data)
	if err != nil {
		return false
	}

	if err := json.Unmarshal(dataBytes, v); err != nil {
		return false
	}

	return true
}

// Delete data by key from cache.
func (c *Cache) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.data, key)
}

// Expired check data by key from cache.
func (c *Cache) Expired(key string) bool {
	c.m.RLock()
	defer c.m.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return true
	}
	return time.Now().After(item.ExpiresAt)
}

// Cleanup all cache.
func (c *Cache) Cleanup() {
	c.m.Lock()
	defer c.m.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
		}
	}
}

// StartCleanup run cleanup all cache.
func (c *Cache) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.Cleanup()
		}
	}()
}
