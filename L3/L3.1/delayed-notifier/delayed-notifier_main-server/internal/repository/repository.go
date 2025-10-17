package repository

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis"
)

type Repository struct {
	*rpRedis.RpRedis
}

func New(rd *pkgRedis.Client) (*Repository, error) {
	return &Repository{
		RpRedis: rpRedis.New(rd),
	}, nil

}
