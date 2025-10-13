package redis

import (
	"context"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Cache — комбинированный кэш: локальный (LRU) + Redis.
type Cache struct {
	client   *Client
	local    *lru.Cache[string, cacheItem]
	mu       sync.RWMutex
	defaultT time.Duration
}

type cacheItem struct {
	value      string
	expiration time.Time
}

// NewCache создаёт новый комбинированный кэш.
func NewCache(client *Client, size int, defaultTTL time.Duration) (*Cache, error) {
	lc, err := lru.New[string, cacheItem](size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		client:   client,
		local:    lc,
		defaultT: defaultTTL,
	}, nil
}

// Set сохраняет значение и в локальный, и в Redis-кэш.
func (c *Cache) Set(ctx context.Context, key, value string, ttl ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := c.defaultT
	if len(ttl) > 0 {
		t = ttl[0]
	}

	// обновляем локальный кэш
	c.local.Add(key, cacheItem{
		value:      value,
		expiration: time.Now().Add(t),
	})

	// записываем в Redis
	return c.client.Set(ctx, key, value, t)
}

// Get возвращает значение из локального кэша или Redis.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	if item, ok := c.local.Get(key); ok {
		if time.Now().Before(item.expiration) {
			c.mu.RUnlock()
			return item.value, nil
		}
	}
	c.mu.RUnlock()

	// Если в локальном нет — читаем из Redis
	val, err := c.client.Get(ctx, key)
	if err != nil || val == "" {
		return "", err
	}

	// Кладём обратно в локальный кэш
	c.mu.Lock()
	c.local.Add(key, cacheItem{
		value:      val,
		expiration: time.Now().Add(c.defaultT),
	})
	c.mu.Unlock()

	return val, nil
}

// Del удаляет значение из обоих уровней.
func (c *Cache) Del(ctx context.Context, key string) error {
	c.mu.Lock()
	c.local.Remove(key)
	c.mu.Unlock()
	return c.client.Del(ctx, key)
}
