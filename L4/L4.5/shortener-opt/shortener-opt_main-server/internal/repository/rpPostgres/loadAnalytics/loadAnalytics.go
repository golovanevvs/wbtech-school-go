package loadAnalytics

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRetry"
)

type RpPostgresLoadAnalytics struct {
	pg *pkgPostgres.Postgres
	rs *pkgRetry.Retry
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgresLoadAnalytics {
	return &RpPostgresLoadAnalytics{
		pg: pg,
		rs: rs,
	}
}

func (rp *RpPostgresLoadAnalytics) LoadAnalytics(ctx context.Context, short string) (totalClicks int, events []model.Analitic, err error) {
	countQuery := `SELECT COUNT(*) FROM analytic WHERE short = $1`
	err = rp.pg.DB.Master.QueryRowContext(ctx, countQuery, short).Scan(&totalClicks)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	eventsQuery := `
		SELECT id, short, user_agent, ip, created_at
		FROM analytic
		WHERE short = $1
		ORDER BY created_at DESC
	`

	rows, err := rp.pg.DB.Master.QueryContext(ctx, eventsQuery, short)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to query analytics events: %w", err)
	}
	defer rows.Close()

	events = make([]model.Analitic, 0)
	for rows.Next() {
		var event model.Analitic
		err := rows.Scan(&event.ID, &event.Short, &event.UserAgent, &event.IP, &event.CreatedAt)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to scan analytics event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return 0, nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return totalClicks, events, nil
}
