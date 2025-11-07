package updateNoticeService

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type IUpdateNoticeRepository interface {
	UpdateNotice(ctx context.Context, notice *model.Notice) (err error)
}

type UpdateNoticeService struct {
	lg *zlog.Zerolog
	rb *pkgRabbitmq.Client
	rp IUpdateNoticeRepository
}

func New(
	parentLg *zlog.Zerolog,
	rp IUpdateNoticeRepository,
) *UpdateNoticeService {
	lg := parentLg.With().Str("component", "AddNoticeService").Logger()
	return &UpdateNoticeService{
		lg: &lg,
		rp: rp,
	}
}

func (sv *UpdateNoticeService) UpdateStatus(ctx context.Context, notice *model.Notice, newStatus model.Status) (err error) {
	lg := sv.lg.With().Str("method", "UpdateNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	notice.Status = newStatus

	lg.Trace().Msgf("%s updating notice status to repository...", pkgConst.OpStart)
	err = sv.rp.UpdateNotice(ctx, notice)
	if err != nil {
		return pkgErrors.Wrap(err, "update notice status to repository")
	}
	lg.Trace().Msgf("%s notice status updated to repository successfully", pkgConst.OpSuccess)

	return nil
}
