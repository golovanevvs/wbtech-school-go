package deleteNoticeService

import (
	"context"
	"errors"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

type IDelRepository interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type IGetService interface {
	GetNotice(ctx context.Context, id int) (notice *model.Notice, err error)
}

type IUpdService interface {
	UpdateStatus(ctx context.Context, notice *model.Notice, newStatus model.Status) (err error)
}

type DeleteNoticeService struct {
	lg    *zlog.Zerolog
	rpDel IDelRepository
	svGet IGetService
	svUpd IUpdService
}

func New(parentLg *zlog.Zerolog, rpDel IDelRepository, svGet IGetService, svUpd IUpdService) *DeleteNoticeService {
	lg := zlog.Logger.With().Str("component", "deleteNoticeService").Logger()
	return &DeleteNoticeService{
		lg:    &lg,
		rpDel: rpDel,
		svGet: svGet,
		svUpd: svUpd,
	}
}

func (sv *DeleteNoticeService) DeleteNotice(ctx context.Context, id int) (err error) {
	lg := sv.lg.With().Str("method", "DeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s deleting notice from repository...", pkgConst.OpStart)
	err = sv.rpDel.DeleteNotice(ctx, id)
	if err != nil {
		return pkgErrors.Wrap(err, "delete notice")
	}
	lg.Trace().Msgf("%s notice deleted from repository successfully", pkgConst.OpSuccess)

	return nil
}

func (sv *DeleteNoticeService) PreDeleteNotice(ctx context.Context, id int) (err error) {
	lg := sv.lg.With().Str("method", "PreDeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Int("notice ID", id).Msgf("%s getting notice from repository...", pkgConst.OpStart)
	notice, err := sv.svGet.GetNotice(ctx, id)
	if err != nil {
		return pkgErrors.Wrap(err, "predelete notice")
	}
	lg.Trace().Int("notice ID", id).Msgf("%s notice got from repository successfully", pkgConst.OpSuccess)

	if notice.Status != model.StatusScheduled {
		return errors.New("failed to delete notice, notice status: " + string(notice.Status))
	}

	lg.Trace().Msgf("%s updating notice status to repository...", pkgConst.OpStart)
	if err := sv.svUpd.UpdateStatus(ctx, notice, model.StatusDeleted); err != nil {
		return pkgErrors.Wrap(err, "predelete notice")
	}
	lg.Trace().Msgf("%s notice status updated from repository successfully", pkgConst.OpSuccess)

	return nil
}
