package repository

import (
	"errors"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis"
	"github.com/wb-go/wbf/zlog"
)

type Repository struct {
	*rpRedis.RpRedis
}

func New(rd *pkgRedis.Client) (*Repository, error) {
	if rd == nil {
		return nil, errors.New("Redis client is nil")
	}
	lg := zlog.Logger.With().Str("layer", "repository").Logger()
	return &Repository{
		RpRedis: rpRedis.New(&lg, rd),
	}, nil

}
