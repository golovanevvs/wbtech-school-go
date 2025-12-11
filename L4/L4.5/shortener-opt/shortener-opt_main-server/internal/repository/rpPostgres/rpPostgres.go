package rpPostgres

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository/rpPostgres/loadAnalytics"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository/rpPostgres/loadOriginalURL"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository/rpPostgres/saveClickEvent"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository/rpPostgres/saveShortURL"
)

type RpPostgres struct {
	*saveShortURL.RpPostgresSaveShortURL
	*loadOriginalURL.RpPostgresLoadOriginalURL
	*saveClickEvent.RpPostgresSaveClickEvent
	*loadAnalytics.RpPostgresLoadAnalytics
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgres {
	return &RpPostgres{
		RpPostgresSaveShortURL:    saveShortURL.New(pg, rs),
		RpPostgresLoadOriginalURL: loadOriginalURL.New(pg, rs),
		RpPostgresLoadAnalytics:   loadAnalytics.New(pg, rs),
		RpPostgresSaveClickEvent:  saveClickEvent.New(pg, rs),
	}
}
