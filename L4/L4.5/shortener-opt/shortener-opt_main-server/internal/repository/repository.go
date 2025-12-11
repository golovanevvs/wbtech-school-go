package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) (*Repository, error) {
	return &Repository{
		RpPostgres: rpPostgres.New(pg, rs),
	}, nil

}
