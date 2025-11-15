package saveShortURL

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRetry"
	"github.com/wb-go/wbf/retry"
)

type RpPostgresSaveShortURL struct {
	pg *pkgPostgres.Postgres
	rs *pkgRetry.Retry
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgresSaveShortURL {
	return &RpPostgresSaveShortURL{
		pg: pg,
		rs: rs,
	}
}

func (rp *RpPostgresSaveShortURL) SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error) {
	query := `
		INSERT INTO short_url (original, short, custom)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var newID int
	row, err := rp.pg.DB.QueryRowWithRetry(ctx, retry.Strategy(*rp.rs), query, shortURL.Original, shortURL.Short, shortURL.Custom)
	if err != nil {
		return 0, fmt.Errorf("failed to save short URL after retries: %w", err)
	}

	err = row.Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to scan ID: %w", err)
	}

	return newID, nil
}
