package deleteNoticePostgres

import (
	"context"
)

type DeleteNoticePostgres struct {
}

func New() *DeleteNoticePostgres {
	return &DeleteNoticePostgres{}
}

func (p *DeleteNoticePostgres) DeleteNotice(ctx context.Context, id int) (err error) {
	return nil
}
