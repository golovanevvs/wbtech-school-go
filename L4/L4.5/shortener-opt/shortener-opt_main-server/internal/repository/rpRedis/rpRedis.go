package rpRedis

import (
	"context"
	"errors"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRedis"
	"github.com/redis/go-redis/v9"
)

type RpRedis struct {
	client *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedis {
	return &RpRedis{
		client: rd,
	}
}

func (rd *RpRedis) GetOriginalURL(ctx context.Context, short string) (string, error) {
	url, err := rd.client.Get(ctx, "url:"+short)
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", pkgErrors.Wrap(err, "get original URL from Redis")
	}
	return url, nil
}
