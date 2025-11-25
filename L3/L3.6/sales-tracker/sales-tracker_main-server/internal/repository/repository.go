package repository

import (
	rpPostgres "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/repository/rpPostgres"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/pkg/pkgPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
}

func New(pg *pkgPostgres.Postgres) *Repository {
	rpPostgres := rpPostgres.New(pg)

	return &Repository{
		RpPostgres: rpPostgres,
	}
}
