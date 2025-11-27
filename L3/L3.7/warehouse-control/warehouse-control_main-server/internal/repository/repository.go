package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
	*rpPostgres.UserRepository
	*rpPostgres.RefreshTokenRepository
}

func New(pg *pkgPostgres.Postgres) (*Repository, error) {
	userRp := rpPostgres.NewUserRepository(pg)
	return &Repository{
		RpPostgres:             rpPostgres.New(pg),
		UserRepository:         userRp,
		RefreshTokenRepository: rpPostgres.NewRefreshTokenRepository(pg),
	}, nil

}
