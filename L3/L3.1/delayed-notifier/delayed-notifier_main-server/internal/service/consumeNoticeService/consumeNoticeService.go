package consumeNoticeService

import (
	"context"
	"encoding/json"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
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
	lg        *zlog.Zerolog
	rb        *pkgRabbitmq.Client
	delNotSv  iDeleteNoticeService
	sendNotSv iSendNoticeService
}

func New(parentLg *zlog.Zerolog, rb *pkgRabbitmq.Client, delNotSv iDeleteNoticeService, sendNotSv iSendNoticeService) *ConsumeNoticeService {
	lg := parentLg.With().Str("component-1", "ConsumeNoticeService").Logger()
	return &ConsumeNoticeService{
		lg:        &lg,
		rb:        rb,
		delNotSv:  delNotSv,
		sendNotSv: sendNotSv,
	}
}

func (sv *ConsumeNoticeService) Consume(ctx context.Context) error {
	lg := sv.lg.With().Str("method", "Consume").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	handler := func(msg amqp.Delivery) {
		sv.handleMessage(ctx, msg)
	}

	if err := sv.rb.ConsumeDLQWithWorkers(ctx, 5, handler); err != nil {
		return pkgErrors.Wrap(err, "consumeDLQ with workers")
	}

	lg.Info().Msgf("%s consumer started", color.BlueString("ℹ️"))

	return nil
}

func (sv *ConsumeNoticeService) handleMessage(ctx context.Context, message amqp.Delivery) {
	lg := sv.lg.With().Str("method", "handleMessage").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	// getting from DLQ, unmarshaling
	var notice model.Notice
	lg.Trace().Msgf("%s unmarshaling message to notice...", color.YellowString("➤"))
	err := json.Unmarshal(message.Body, &notice)
	if err != nil {
		sv.lg.Error().Err(err).Msgf("%s failed to unmarshal message", color.RedString("❌"))
		return
	}
	lg.Trace().Int("notice ID", notice.ID).Str("status", string(notice.Status)).Msgf("%s data unmarshaled successfully", color.GreenString("✔"))

	lg.Trace().Int("notice ID", notice.ID).Msgf("%s acknowledging message...", color.YellowString("➤"))
	if err := sv.rb.Ack(message); err != nil {
		sv.lg.Error().Err(err).Int("notice ID", notice.ID).Msgf("%s failed to ack message", color.RedString("❌"))
	}
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s message acknowledged successfully", color.GreenString("✔"))

	// cheking status, setting new status and sending notice
	if notice.Status != model.StatusDeleted {
		notice.Status = model.StatusPending
		lg.Trace().Int("notice ID", notice.ID).Msgf("%s sending message...", color.YellowString("➤"))
		sv.sendNotSv.SendNotice(ctx, notice)
		lg.Trace().Int("notice ID", notice.ID).Msgf("%s message sending completed", color.GreenString("✔"))
		notice.Status = model.StatusSent
	}

	// deleting notice from repository
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s deleting message from repository...", color.YellowString("➤"))
	if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
		sv.lg.Error().Err(err).Int("notice ID", notice.ID).Msgf("%s failed to delete notice from repository", color.RedString("❌"))
		return
	}
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s message deleted from repository", color.GreenString("✔"))
}
