package rediscache

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
	log    *zerolog.Logger
}
