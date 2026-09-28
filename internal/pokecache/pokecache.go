package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Adds a value to the cache.
func (c *Cache) Add(key string, v []byte) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: v}
	c.mu.Unlock()
}

// Gets and returns a value from the cache, and whether or not it exists.
func (c *Cache) Get(key string) (v []byte, exists bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// Loops over the cache as long as it exists and removes entries older than `interval`.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		c.mu.Lock()

		for key, entry := range c.entries {
			if t.Sub(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()
	}
}

// Creates a new Cache, and returns a pointer to it.
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{entries: make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return cache
}
