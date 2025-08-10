package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
	log    *zerolog.Logger
}

// ErrRedisNil is returned when a key is not found in Redis.
const ErrRedisNil = redis.Nil

func New(redisCfg *config.RedisCache, logger *zerolog.Logger) (*RedisCache, error) {
	log := logger.With().Str("component", "redis").Logger()

	rd := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rd.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("failed to connect to Redis")
		return nil, err
	}

	log.Info().
		Str("addr", redisCfg.Addr).
		Dur("ttl", redisCfg.TTL).
		Msg("connected to Redis")

	return &RedisCache{
		client: rd,
		ttl:    redisCfg.TTL,
		log:    &log,
	}, nil
}

func (rd *RedisCache) Set(ctx context.Context, key string, value any) error {
	log := rd.log.With().Str("key", key).Logger()

	data, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Msg("failed to encode value for Redis")
		return err
	}

	if err := rd.client.Set(ctx, key, data, rd.ttl).Err(); err != nil {
		log.Error().Err(err).Msg("failed to set key in Redis")
		return err
	}

	log.Debug().Msg("value set in Redis")

	return nil
}

func (rd *RedisCache) SetOrder(ctx context.Context, key string, order *model.Order) error {
	return rd.Set(ctx, key, order)
}

func (rd *RedisCache) Get(ctx context.Context, key string, dest any) error {
	log := rd.log.With().Str("key", key).Logger()

	data, err := rd.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Debug().Msg("key not found in Redis")
			return ErrRedisNil
		}

		log.Error().Err(err).Msg("failed to get key from Redis")

		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		log.Error().Err(err).Msg("failed to decode Redis value")
		return err
	}

	log.Debug().Msg("value retrieved from Redis")

	return nil
}

func (rd *RedisCache) GetOrder(ctx context.Context, key string) (*model.Order, error) {
	var order model.Order
	err := rd.Get(ctx, key, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (rd *RedisCache) Delete(ctx context.Context, key string) error {
	log := rd.log.With().Str("key", key).Logger()

	if err := rd.client.Del(ctx, key).Err(); err != nil {
		log.Error().Err(err).Msg("failed to delete key from Redis")
		return err
	}

	log.Debug().Msg("key deleted from Redis")

	return nil
}

func (rd *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	log := rd.log.With().Str("key", key).Logger()

	count, err := rd.client.Exists(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to check key existence")
		return false, err
	}

	return count > 0, nil
}

func (rd *RedisCache) Close() error {
	err := rd.client.Close()
	if err != nil {
		rd.log.Error().Err(err).Msg("failed to close Redis client")
		return err
	}

	rd.log.Info().Msg("Redis client closed")

	return nil
}
