package consumeNoticeService

import (
	"context"
	"encoding/json"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/zlog"
)

type iDeleteNoticeService interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type iSendNoticeService interface {
	SendNotice(ctx context.Context, notice model.Notice)
}

type ConsumeNoticeService struct {
	lg        zlog.Zerolog
	rb        *pkgRabbitmq.Client
	delNotSv  iDeleteNoticeService
	sendNotSv iSendNoticeService
}

func New(rb *pkgRabbitmq.Client, delNotSv iDeleteNoticeService, sendNotSv iSendNoticeService) *ConsumeNoticeService {
	lg := zlog.Logger.With().Str("component", "service-consumeNoticeService").Logger()
	return &ConsumeNoticeService{
		lg:        lg,
		rb:        rb,
		delNotSv:  delNotSv,
		sendNotSv: sendNotSv,
	}
}

func (sv *ConsumeNoticeService) Consume(ctx context.Context) error {
	sv.lg.Debug().Msg("----- consumer starting...")
	handler := func(msg amqp.Delivery) {
		sv.handleMessage(ctx, msg)
	}

	if err := sv.rb.ConsumeDLQWithWorkers(ctx, 5, handler); err != nil {
		sv.lg.Error().Err(err).Msg("failed to consume DLQ with workers")
		return err
	}

	sv.lg.Info().Msg("----- consumer started")

	return nil
}

func (sv *ConsumeNoticeService) handleMessage(ctx context.Context, message amqp.Delivery) {
	sv.lg.Trace().Msg("--- consume handler started")
	defer sv.lg.Trace().Msg("--- consume handler stopped")

	// getting from DLQ, unmarshaling
	var notice model.Notice
	err := json.Unmarshal(message.Body, &notice)
	if err != nil {
		sv.lg.Error().Err(err).Msg("failed to unmarshal message")
		return
	}
	sv.lg.Debug().Str("message", notice.Message).Msg("received message")
	if err := sv.rb.Ack(message); err != nil {
		sv.lg.Error().Err(err).Msg("failed to ack message")
	}

	// sending
	sv.sendNotSv.SendNotice(ctx, notice)

	// deleting from repository
	if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
		sv.lg.Error().Err(err).Msg("failed to delete notice")
		return
	}
}
