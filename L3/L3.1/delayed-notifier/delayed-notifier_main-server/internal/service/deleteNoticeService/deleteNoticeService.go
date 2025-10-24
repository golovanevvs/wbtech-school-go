package deleteNoticeService

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type DeleteNoticeService struct {
	lg *zlog.Zerolog
	rp IRepository
}

func New(rp IRepository) *DeleteNoticeService {
	lg := zlog.Logger.With().Str("component", "deleteNoticeService").Logger()
	return &DeleteNoticeService{
		lg: &lg,
		rp: rp,
	}
}

func (sv *DeleteNoticeService) DeleteNotice(ctx context.Context, id int) (err error) {
	lg := sv.lg.With().Str("method", "DeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s deleting notice from repository...", pkgConst.OpStart)
	err = sv.rp.DeleteNotice(ctx, id)
	if err != nil {
		return pkgErrors.Wrap(err, "delete notice")
	}
	lg.Trace().Msgf("%s notice deleted from repository successfully", pkgConst.OpSuccess)

	return nil
}
