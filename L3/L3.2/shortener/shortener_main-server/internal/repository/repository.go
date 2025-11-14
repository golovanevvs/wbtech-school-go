package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository/rpPostgres"
	"github.com/wb-go/wbf/zlog"
)

type Repository struct {
	*rpPostgres.RpPostgres
}

func New() (*Repository, error) {
	lg := zlog.Logger.With().Str("layer", "repository").Logger()
	return &Repository{
		RpPostgres: rpPostgres.New(&lg),
	}, nil

}
