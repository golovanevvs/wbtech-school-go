package loadOriginalURL

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
)

type RpPostgresLoadOriginalURL struct {
	pg *pkgPostgres.Postgres
	rs *pkgRetry.Retry
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgresLoadOriginalURL {
	return &RpPostgresLoadOriginalURL{
		pg: pg,
		rs: rs,
	}
}

func (rp *RpPostgresLoadOriginalURL) LoadOriginalURL(ctx context.Context, short string) (original string, err error) {
	query := `
		SELECT original
		FROM short_url
		WHERE short = $1
	`

	var originalURL string
	row := rp.pg.DB.Master.QueryRowContext(ctx, query, short)
	err = row.Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("URL not found for short code: %s", short)
		}
		return "", fmt.Errorf("failed to query original URL: %w", err)
	}

	return originalURL, nil
}
