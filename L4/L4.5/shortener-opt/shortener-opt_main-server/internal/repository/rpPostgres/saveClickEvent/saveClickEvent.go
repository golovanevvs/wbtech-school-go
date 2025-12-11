package saveClickEvent

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
	"github.com/wb-go/wbf/retry"
)

type RpPostgresSaveClickEvent struct {
	pg *pkgPostgres.Postgres
	rs *pkgRetry.Retry
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgresSaveClickEvent {
	return &RpPostgresSaveClickEvent{
		pg: pg,
		rs: rs,
	}
}

func (rp *RpPostgresSaveClickEvent) SaveClickEvent(ctx context.Context, event model.Analitic) error {
	query := `
		INSERT INTO analytic (short, user_agent, ip)
		VALUES ($1, $2, $3)
	`

	_, err := rp.pg.DB.ExecWithRetry(ctx, retry.Strategy(*rp.rs), query, event.Short, event.UserAgent, event.IP)
	if err != nil {
		return fmt.Errorf("failed to save click event: %w", err)
	}

	return nil
}
