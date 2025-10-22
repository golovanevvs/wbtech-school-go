package addNoticeService

import (
	"context"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type ISaveNoticeRepository interface {
	SaveNotice(ctx context.Context, notice model.Notice) (id int, err error)
}

type IDeleteNoticeService interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type AddNoticeService struct {
	lg       *zlog.Zerolog
	rb       *pkgRabbitmq.Client
	delNotSv IDeleteNoticeService
	rp       ISaveNoticeRepository
}

func New(
	parentLg *zlog.Zerolog,
	rb *pkgRabbitmq.Client,
	delNotSv IDeleteNoticeService,
	rp ISaveNoticeRepository,
) *AddNoticeService {
	lg := parentLg.With().Str("component", "AddNoticeService").Logger()
	return &AddNoticeService{
		lg:       &lg,
		rb:       rb,
		delNotSv: delNotSv,
		rp:       rp,
	}
}

func (sv *AddNoticeService) AddNotice(ctx context.Context, reqNotice model.ReqNotice) (id int, err error) {
	lg := sv.lg.With().Str("method", "AddNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	createdAt := time.Now()
	sentAt := reqNotice.SentAt
	ttl := sentAt.Sub(createdAt)
	notice := model.Notice{
		UserID:    reqNotice.UserID,
		Message:   reqNotice.Message,
		Channels:  reqNotice.Channels,
		CreatedAt: createdAt,
		SentAt:    sentAt,
		Status:    model.StatusScheduled,
	}

	lg.Trace().Msgf("%s saving notice to repository...", pkgConst.OpStart)
	id, err = sv.rp.SaveNotice(ctx, notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save notice to repository")
	}
	lg.Trace().Msgf("%s notice saved to repository successfully", pkgConst.OpSuccess)

	notice.ID = id

	lg.Trace().Int("notice ID", notice.ID).Msgf("%s publishing notice with TTL to message broker...", pkgConst.OpStart)
	if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
		lg.Debug().Err(err).Msgf("%s failed to publish struct with TTL to RabbitMQ", pkgConst.Error)
		if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
			lg.Debug().Err(err).Int("notice ID", notice.ID).Msgf("%s failed deleted notice from Redis", pkgConst.Error)
		}
		return 0, pkgErrors.Wrap(err, "error publish struct with TTL to RabbitMQ")
	}
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s notice with TTL published to message broker successfully", pkgConst.OpSuccess)

	return id, nil
}
