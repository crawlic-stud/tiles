package cache

import (
	"sync"
	"time"
)

// Item represents a cached value with optional expiration.
type Item[V any] struct {
	Value      V
	Expiration int64
}

// Cache is a thread-safe generic in-memory cache.
type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]Item[V]
	ttl   time.Duration
}

// New creates a new cache instance.
func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]Item[V]),
		ttl:   ttl,
	}
}

// Set stores a value in the cache.
// ttl <= 0 means no expiration.
func (c *Cache[K, V]) Set(key K, value V) {
	var exp int64

	if c.ttl > 0 {
		exp = time.Now().Add(c.ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item[V]{
		Value:      value,
		Expiration: exp,
	}
}

// Get returns a value from the cache.
// bool=false means key not found or expired.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	var zero V
	if !ok {
		return zero, false
	}

	// Check expiration
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.Delete(key)
		return zero, false
	}

	return item.Value, true
}

// Delete removes a key from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Exists checks whether a key exists and is not expired.
func (c *Cache[K, V]) Exists(key K) bool {
	_, ok := c.Get(key)
	return ok
}

// Clear removes all items from the cache.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]Item[V])
}

// Size returns the number of items in the cache.
func (c *Cache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
