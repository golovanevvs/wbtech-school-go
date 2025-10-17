package addNoticeService

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
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
	lg       zlog.Zerolog
	rp       ISaveNoticeRepository
	rb       *pkgRabbitmq.Client
	delNotSv IDeleteNoticeService
}

func New(rp ISaveNoticeRepository, rb *pkgRabbitmq.Client, delNotSv IDeleteNoticeService) *AddNoticeService {
	lg := zlog.Logger.With().Str("component", "service-addNoticeService").Logger()
	return &AddNoticeService{
		lg:       lg,
		rp:       rp,
		rb:       rb,
		delNotSv: delNotSv,
	}
}

func (sv *AddNoticeService) AddNotice(ctx context.Context, reqNotice model.ReqNotice) (id int, err error) {

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

	sv.lg.Trace().Msgf("%s saving notice to repository...", color.YellowString("➤"))
	id, err = sv.rp.SaveNotice(ctx, notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save notice to repository")
	}
	sv.lg.Trace().Msgf("%s notice saved to repository successfully", color.GreenString("✔"))

	notice.ID = id

	sv.lg.Trace().Int("notice ID", notice.ID).Msgf("%s publishing notice with TTL to message broker...", color.YellowString("➤"))
	if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
		sv.lg.Error().Err(err).Msg("error publish struct with TTL to RabbitMQ")
		if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
			sv.lg.Trace().Err(err).Int("notice ID", notice.ID).Msg("failed deleted notice from Redis")
		}
		return 0, fmt.Errorf("error publish struct with TTL to RabbitMQ")
	}
	sv.lg.Trace().Int("notice ID", notice.ID).Msgf("%s notice with TTL published to message broker successfully", color.GreenString("✔"))

	return id, nil
}
