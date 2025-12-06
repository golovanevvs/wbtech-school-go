package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
}

func New(pg *pkgPostgres.Postgres) (*Repository, error) {
	return &Repository{
		RpPostgres: rpPostgres.New(pg),
	}, nil

}
