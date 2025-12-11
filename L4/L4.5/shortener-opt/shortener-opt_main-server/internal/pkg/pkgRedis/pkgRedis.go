package pkgRedis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgErrors"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

func New(cfg *Config) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{rdb: rdb, ttl: cfg.TTL}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	t := c.ttl
	if len(ttl) > 0 {
		t = ttl[0]
	}
	return c.rdb.Set(ctx, key, value, t).Err()
}

func (c *Client) SetWithID(ctx context.Context, prefix string, value interface{}, ttl ...time.Duration) (string, error) {
	id, err := c.rdb.Incr(ctx, prefix+":next_id").Result()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s:%d", prefix, id)
	t := c.ttl
	if len(ttl) > 0 {
		t = ttl[0]
	}

	if err := c.rdb.Set(ctx, key, value, t).Err(); err != nil {
		return "", err
	}

	return key, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", pkgErrors.Wrap(pkgErrors.ErrNoticeNotFound, "get data from Redis")
	}
	if err != nil {
		return "", pkgErrors.Wrap(err, "get data from Redis")
	}
	return val, nil
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	deleted, err := c.rdb.Del(ctx, keys...).Result()
	if err != nil {
		return pkgErrors.Wrap(err, "delete data from Redis")
	}
	if deleted == 0 {
		return pkgErrors.Wrap(pkgErrors.ErrNoticeNotFound, "delete data from Redis")

	}
	return nil
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, key, ttl).Err()
}

func (c *Client) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return c.rdb.HSet(ctx, key, values).Err()
}

func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := c.rdb.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}
