package saveShortURL

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgPostgres"
	"github.com/wb-go/wbf/zlog"
)

type RpPostgresSaveShortURL struct {
	lg *zlog.Zerolog
	pg *pkgPostgres.Postgres
}

func New(parentLg *zlog.Zerolog, pg *pkgPostgres.Postgres) *RpPostgresSaveShortURL {
	lg := parentLg.With().Str("component", "RpRedisDeleteNotice").Logger()
	return &RpPostgresSaveShortURL{
		lg: &lg,
		pg: pg,
	}
}

func (rp *RpPostgresSaveShortURL) SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error) {
	return 100, nil
}
