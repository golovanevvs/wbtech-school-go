package getNoticeService

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	LoadNotice(ctx context.Context, id int) (notice *model.Notice, err error)
}

type GetNoticeService struct {
	lg *zlog.Zerolog
	rp IRepository
}

func New(parentLg *zlog.Zerolog, rp IRepository) *GetNoticeService {
	lg := parentLg.With().Str("component", "GetNoticeService").Logger()
	return &GetNoticeService{
		lg: &lg,
		rp: rp,
	}
}

func (sv *GetNoticeService) GetNotice(ctx context.Context, id int) (notice *model.Notice, err error) {
	lg := sv.lg.With().Str("method", "GetNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s get notice from repository...", pkgConst.OpStart)
	notice, err = sv.rp.LoadNotice(ctx, id)
	if err != nil {
		return nil, pkgErrors.Wrapf(err, "get notice from repository, notice ID: %d", id)
	}
	lg.Trace().Msgf("%s notice got from repository successfully", pkgConst.OpSuccess)

	return notice, nil
}
