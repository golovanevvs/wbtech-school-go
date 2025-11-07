package consumeNoticeService

import (
	"context"
	"encoding/json"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
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

type iGetNoticeService interface {
	GetNotice(ctx context.Context, id int) (notice *model.Notice, err error)
}

type iUpdateNoticeService interface {
	UpdateStatus(ctx context.Context, notice *model.Notice, newStatus model.Status) (err error)
}

type ConsumeNoticeService struct {
	lg        *zlog.Zerolog
	rb        *pkgRabbitmq.Client
	delNotSv  iDeleteNoticeService
	sendNotSv iSendNoticeService
	getNotSv  iGetNoticeService
	updNotSv  iUpdateNoticeService
}

func New(
	parentLg *zlog.Zerolog,
	rb *pkgRabbitmq.Client,
	delNotSv iDeleteNoticeService,
	sendNotSv iSendNoticeService,
	getNotSv iGetNoticeService,
	updNotSv iUpdateNoticeService,
) *ConsumeNoticeService {
	lg := parentLg.With().Str("component", "ConsumeNoticeService").Logger()
	return &ConsumeNoticeService{
		lg:        &lg,
		rb:        rb,
		delNotSv:  delNotSv,
		sendNotSv: sendNotSv,
		getNotSv:  getNotSv,
		updNotSv:  updNotSv,
	}
}

func (sv *ConsumeNoticeService) Consume(ctx context.Context) error {
	lg := sv.lg.With().Str("method", "Consume").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Info().Msgf("%s consumer starting...", pkgConst.Starting)

	handler := func(msg amqp.Delivery) {
		sv.handleMessage(ctx, msg)
	}

	if err := sv.rb.ConsumeDLQWithWorkers(ctx, 5, handler); err != nil {
		return pkgErrors.Wrap(err, "consumeDLQ with workers")
	}

	lg.Info().Msgf("%s consumer started successfully", pkgConst.Finished)

	return nil
}

func (sv *ConsumeNoticeService) handleMessage(ctx context.Context, message amqp.Delivery) {
	lg := sv.lg.With().Str("method", "handleMessage").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	// getting from DLQ, unmarshaling
	var notice model.Notice
	lg.Trace().Msgf("%s unmarshaling message to notice...", pkgConst.OpStart)
	err := json.Unmarshal(message.Body, &notice)
	if err != nil {
		sv.lg.Error().Err(err).Msgf("%s failed to unmarshal message", pkgConst.Error)
		return
	}
	lg.Trace().Int("notice ID", notice.ID).Str("status", string(notice.Status)).Msgf("%s data unmarshaled successfully", pkgConst.OpSuccess)

	lg.Trace().Int("notice ID", notice.ID).Msgf("%s acknowledging message...", pkgConst.OpStart)
	if err := sv.rb.Ack(message); err != nil {
		sv.lg.Error().Err(err).Int("notice ID", notice.ID).Msgf("%s failed to ack message", pkgConst.Error)
	}
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s message acknowledged successfully", pkgConst.OpSuccess)

	// cheking status, setting new status and sending notice
	updateNotice, err := sv.getNotSv.GetNotice(ctx, notice.ID)
	if err != nil {
		lg.Error().Err(err).Int("notice ID", notice.ID).Msg("failed to update notice status")
		return
	}
	if updateNotice.Status != model.StatusDeleted {
		if err := sv.updNotSv.UpdateStatus(ctx, &notice, model.StatusPending); err != nil {
			lg.Error().Err(err).Int("notice ID", notice.ID).Msg("failed to update notice status")
			return
		}
		lg.Trace().Int("notice ID", notice.ID).Msgf("%s sending message...", pkgConst.OpStart)
		sv.sendNotSv.SendNotice(ctx, notice)
		lg.Trace().Int("notice ID", notice.ID).Msgf("%s message sending completed", pkgConst.OpSuccess)
		notice.Status = model.StatusSent
	}

	// deleting notice from repository
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s deleting message from repository...", pkgConst.OpStart)
	if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
		sv.lg.Error().Err(err).Int("notice ID", notice.ID).Msgf("%s failed to delete notice from repository", pkgConst.Error)
		return
	}
	lg.Trace().Int("notice ID", notice.ID).Msgf("%s message deleted from repository", pkgConst.OpSuccess)
}
