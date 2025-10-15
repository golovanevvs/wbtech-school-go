package consumeNoticeService

import (
	"context"
	"encoding/json"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/zlog"
)

type iRepository interface {
	LoadTelName(ctx context.Context, username string) (chatID int, err error)
}

type ConsumeNoticeService struct {
	lg zlog.Zerolog
	rb *pkgRabbitmq.Client
	tg *pkgTelegram.Client
	rd *pkgRedis.Client
}

func New(rb *pkgRabbitmq.Client, tg *pkgTelegram.Client, rd *pkgRedis.Client) *ConsumeNoticeService {
	lg := zlog.Logger.With().Str("component", "service-consumeNoticeService").Logger()
	return &ConsumeNoticeService{
		lg: lg,
		rb: rb,
		tg: tg,
		rp: rp,
	}
}

func (sv *ConsumeNoticeService) Consume(ctx context.Context) error {
	sv.lg.Debug().Msg("----- consumer starting...")

	handler := func(message amqp.Delivery) {
		sv.lg.Trace().Msg("--- consume handler started")
		defer sv.lg.Trace().Msg("--- consume handler stopped")

		var notice model.Notice
		err := json.Unmarshal(message.Body, &notice)
		if err != nil {
			sv.lg.Error().Err(err).Msg("failed to unmarshal message")
			return
		}
		sv.lg.Debug().Str("message", notice.Message).Msg("received message")

		for _, ch := range notice.Channels {
			switch ch.Type {
			case model.ChannelTelegram:

				chatID, err := sv.rp.LoadTelName(ctx, ch.Value)
				if err != nil {
					sv.lg.Error().Err(err).Msg("failed to load telegram chat id")
					break
				}

				sv.tg.SendTo(int64(chatID), notice.Message)

				if err := sv.rb.Ack(message); err != nil {
					sv.lg.Error().Err(err).Msg("failed to ack message")
				}
			}
		}
	}

	if err := sv.rb.ConsumeDLQWithWorkers(ctx, 5, handler); err != nil {
		sv.lg.Error().Err(err).Msg("failed to consume DLQ with workers")
		return err
	}

	sv.lg.Info().Msg("----- consumer started")

	return nil
}
