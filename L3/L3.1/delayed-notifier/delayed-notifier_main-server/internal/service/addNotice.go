package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/model"
)

func (sv *Service) AddNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	sv.rp.Postgres.Master.QueryRowContext(ctx, `
	
		INSERT
	
	`)
	return
}
