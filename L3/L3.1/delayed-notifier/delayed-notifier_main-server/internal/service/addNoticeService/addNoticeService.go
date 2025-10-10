package addNoticeService

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	AddNotice(ctx context.Context, notice model.Notice) (id int, err error)
}

type IRabbitMQ interface {
	PublishStructWithTTL(data any, ttl time.Duration) (err error)
}

type AddNoticeService struct {
	lg zlog.Zerolog
	rp IRepository
	rb IRabbitMQ
}

func New(rp IRepository, rb IRabbitMQ) *AddNoticeService {
	lg := zlog.Logger.With().Str("component", "service-addNoticeService").Logger()
	return &AddNoticeService{
		lg: lg,
		rp: rp,
		rb: rb,
	}
}

func (sv *AddNoticeService) AddNotice(ctx context.Context, reqNotice model.Notice) (id int, err error) {
	sv.lg.Trace().Msg("run AddNotice")
	id = 1
	createdAt := time.Now()
	sentAt := reqNotice.SentAt
	ttl := sentAt.Sub(createdAt)
	notice := model.Notice{
		ID:        id,
		UserID:    reqNotice.UserID,
		Message:   reqNotice.Message,
		Channels:  reqNotice.Channels,
		CreatedAt: createdAt,
		SentAt:    sentAt,
		Status:    model.StatusScheduled,
	}

	sv.lg.Trace().Msg("publish struct with TTL to RabbitMQ")
	if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
		sv.lg.Error().Err(err).Msg("error publish struct with TTL to RabbitMQ")
		return 0, fmt.Errorf("error publish struct with TTL to RabbitMQ")
	}
	sv.lg.Trace().Msg("struct with TTL published to RabbitMQ successfully")

	return id, nil
}
