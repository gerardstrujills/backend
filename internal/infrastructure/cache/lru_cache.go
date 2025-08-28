package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gerardstrujills/backend/internal/domain/entities"
	lru "github.com/hashicorp/golang-lru/v2"
)

type LRUCache struct {
	cache *lru.Cache[string, *entities.CacheItem]
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewLRUCache(size int, ttl time.Duration) (*LRUCache, error) {
	cache, err := lru.New[string, *entities.CacheItem](size)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el cache LRU: %w", err)
	}

	lruCache := &LRUCache{
		cache: cache,
		ttl:   ttl,
	}

	// Iniciar limpieza periódica de elementos expirados
	go lruCache.cleanupExpired()

	return lruCache, nil
}

func (c *LRUCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	// Verificar si el item ha expirado
	if item.IsExpired() {
		c.mutex.RUnlock()
		c.mutex.Lock()
		c.cache.Remove(key)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, false
	}

	return item.Data, true
}

func (c *LRUCache) Set(ctx context.Context, key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item := &entities.CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	c.cache.Add(key, item)
	return nil
}

func (c *LRUCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache.Remove(key)
	return nil
}

func (c *LRUCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache.Purge()
	return nil
}

// Limpia periódicamente los elementos expirados
func (c *LRUCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		keys := c.cache.Keys()
		for _, key := range keys {
			if item, found := c.cache.Peek(key); found && item.IsExpired() {
				c.cache.Remove(key)
			}
		}
		c.mutex.Unlock()
	}
}
