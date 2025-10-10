package postgres

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/postgres/deleteNoticePostgres"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/postgres/saveNoticePostgres"
)

type Postgres struct {
	SaveNoticePostgres   *saveNoticePostgres.SaveNoticePostgres
	DeleteNoticePostgres *deleteNoticePostgres.DeleteNoticePostgres
}

func New() *Postgres {
	return &Postgres{
		SaveNoticePostgres:   saveNoticePostgres.New(),
		DeleteNoticePostgres: deleteNoticePostgres.New(),
	}
}
