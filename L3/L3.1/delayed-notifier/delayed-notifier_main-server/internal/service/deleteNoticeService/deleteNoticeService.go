package deleteNoticeService

import (
	"context"

	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type DeleteNoticeService struct {
	lg zlog.Zerolog
	rp IRepository
}

func New(rp IRepository) *DeleteNoticeService {
	lg := zlog.Logger.With().Str("component", "service-deleteNoticeService").Logger()
	return &DeleteNoticeService{
		lg: lg,
		rp: rp,
	}
}

func (sv *DeleteNoticeService) DeleteNotice(ctx context.Context, id int) (err error) {
	sv.lg.Trace().Msg("run AddNotice")
	return nil
}
