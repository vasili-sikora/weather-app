package cli

import (
	"context"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (float32, bool, error)
	Set(ctx context.Context, key string, value float32, ttl time.Duration) error
}

type memoryCacheItem struct {
	value     float32
	expiresAt time.Time
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]memoryCacheItem
}

func newMemoryCache() *memoryCache {
	return &memoryCache{
		items: make(map[string]memoryCacheItem),
	}
}

func (c *memoryCache) Get(_ context.Context, key string) (float32, bool, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return 0, false, nil
	}

	if time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return 0, false, nil
	}

	return item.value, true, nil
}

func (c *memoryCache) Set(_ context.Context, key string, value float32, ttl time.Duration) error {
	c.mu.Lock()
	c.items[key] = memoryCacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()

	return nil
}
