package addNoticeService

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	SaveNotice(ctx context.Context, notice model.Notice) (id int, err error)
}

type IService interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type AddNoticeService struct {
	lg       zlog.Zerolog
	rp       IRepository
	rb       *pkgRabbitmq.Client
	delNotSv IService
}

func New(rp IRepository, rb *pkgRabbitmq.Client, delNotSv IService) *AddNoticeService {
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

	sv.lg.Trace().Int("notice ID", notice.ID).Msg("saving notice to Redis...")
	id, err = sv.rp.SaveNotice(ctx, notice)
	if err != nil {
		sv.lg.Error().Err(err).Msg("error save notice")
		return 0, err
	}
	sv.lg.Trace().Msg("notice saved to Redis successfully")

	notice.ID = id

	sv.lg.Trace().Msg("publishing notice with TTL to RabbitMQ...")
	if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
		sv.lg.Error().Err(err).Msg("error publish struct with TTL to RabbitMQ")
		if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
			sv.lg.Trace().Err(err).Msg("failed deleted notice from Redis")
		}
		return 0, fmt.Errorf("error publish struct with TTL to RabbitMQ")
	}
	sv.lg.Trace().Msg("notice with TTL published to RabbitMQ successfully")

	return id, nil
}
