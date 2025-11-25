package repository

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/pkg/pkgPostgres"

type Repository struct {
	*rpPostgres.RpPostgres
}

func New(pg *pkgPostgres.Postgres) (*Repository, error) {
	userRp := rpPostgres.NewUserRepository(pg)
	return &Repository{
		RpPostgres: rpPostgres.New(pg),
	}, nil

}
