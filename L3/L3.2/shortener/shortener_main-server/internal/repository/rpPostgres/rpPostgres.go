package rpPostgres

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository/rpPostgres/saveShortURL"
	"github.com/wb-go/wbf/zlog"
)

type RpPostgres struct {
	*saveShortURL.RpPostgresSaveShortURL
}

func New(parentLg *zlog.Zerolog, pg *pkgPostgres.Postgres) *RpPostgres {
	lg := parentLg.With().Str("component", "RpPostgres").Logger()
	return &RpPostgres{
		RpPostgresSaveShortURL: saveShortURL.New(&lg, pg),
	}
}
