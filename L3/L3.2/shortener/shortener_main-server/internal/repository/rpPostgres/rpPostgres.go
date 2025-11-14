package rpPostgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository/rpPostgres/saveShortURL"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wb-go/wbf/zlog"
)

type RpPostgres struct {
	db *sqlx.DB
	*saveShortURL.RpPostgresSaveShortURL
}

func New(ctx context.Context, parentLg *zlog.Zerolog, cgf *Config) *RpPostgres {
	lg := parentLg.With().Str("component", "RpPostgres").Logger()
	lg.Debug().Msgf("%s connecting to PostreSQL", pkgConst.AppStart)

	// db, err := sqlx.ConnectContext(ctx, "pgx", config.)

	return &RpPostgres{
		RpPostgresSaveShortURL: saveShortURL.New(&lg),
	}
}
