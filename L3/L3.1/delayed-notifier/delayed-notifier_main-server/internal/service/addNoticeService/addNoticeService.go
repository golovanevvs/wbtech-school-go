package addNoticeService

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	SaveNotice(ctx context.Context, notice model.Notice) (id int, err error)
}

type AddNoticeService struct {
	lg zlog.Zerolog
	rp IRepository
	rb *rabbitmq.Client
}

func New(rp IRepository, rb *rabbitmq.Client) *AddNoticeService {
	lg := zlog.Logger.With().Str("component", "service-addNoticeService").Logger()
	return &AddNoticeService{
		lg: lg,
		rp: rp,
		rb: rb,
	}
}

func (sv *AddNoticeService) AddNotice(ctx context.Context, reqNotice model.ReqNotice) (id int, err error) {
	sv.lg.Trace().Msg("AddNotice run...")
	defer sv.lg.Trace().Msg("AddNotice stopped")
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

	sv.lg.Trace().Msg("save notice to Redis")
	id, err = sv.rp.SaveNotice(ctx, notice)
	if err != nil {
		sv.lg.Error().Err(err).Msg("error save notice")
		return 0, err
	}
	sv.lg.Trace().Msg("notice saved to Redis successfully")

	notice.ID = id

	sv.lg.Trace().Msg("publish notice with TTL to RabbitMQ")
	if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
		// удалить из Redis
		sv.lg.Error().Err(err).Msg("error publish struct with TTL to RabbitMQ")
		return 0, fmt.Errorf("error publish struct with TTL to RabbitMQ")
	}
	sv.lg.Trace().Msg("notice with TTL published to RabbitMQ successfully")

	return id, nil
}
