package saveNoticePostgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
)

type SaveNoticePostgres struct {
}

func New() *SaveNoticePostgres {
	return &SaveNoticePostgres{}
}

func (p SaveNoticePostgres) SaveNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	return 0, nil
}
