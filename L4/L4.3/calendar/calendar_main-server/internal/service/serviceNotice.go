package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

// ISendNoticeRp interface for send notice repository
type ISendNoticeRp interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
}

// NoticeService handles sending notifications
type NoticeService struct {
	lg *zlog.Zerolog
	rs *pkgRetry.Retry
	tg *pkgTelegram.Client
	em *pkgEmail.Client
	rp ISendNoticeRp
}

// NewSendNoticeService creates a new SendNoticeService
func NewNoticeService(
	parentLg *zlog.Zerolog,
	rs *pkgRetry.Retry,
	tg *pkgTelegram.Client,
	em *pkgEmail.Client,
	rp ISendNoticeRp,
) *NoticeService {
	lg := parentLg.With().Str("component", "SendNoticeService").Logger()
	return &NoticeService{
		lg: &lg,
		rs: rs,
		tg: tg,
		em: em,
		rp: rp,
	}
}

// SendNotice sends a notification to the user
func (sv *NoticeService) SendNotice(ctx context.Context, notice model.Notice) {
	lg := sv.lg.With().Str("method", "SendNotice").Logger()

	user, err := sv.rp.GetByID(ctx, notice.UserID)
	if err != nil {
		lg.Error().Err(err).Int("user_id", notice.UserID).Msg("Failed to get user for notice")
		return
	}

	wg := sync.WaitGroup{}

	if user.TelegramChatID != nil && notice.Channels.Telegram {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendNoticeToTelegram(ctx, *user.TelegramChatID, notice.Message)
		}()
	}

	if notice.Channels.Email {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendNoticeToEmail(ctx, user.Email, notice.Message)
		}()
	}

	wg.Wait()
}

// SendNoticeToTelegram sends a notification to Telegram
func (sv *NoticeService) SendNoticeToTelegram(ctx context.Context, chatID int64, message string) error {
	lg := sv.lg.With().Str("method", "SendNoticeToTelegram").Logger()

	fn := func() error {
		err := sv.tg.SendTo(chatID, message)
		if err != nil {
			lg.Warn().Err(err).Int64("chat_id", chatID).Msg("Failed to send notice to telegram")
			return err
		}
		return nil
	}

	if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
		return fmt.Errorf("send notice to telegram after all attempts: %w", err)
	}

	lg.Debug().Int64("chat_id", chatID).Msg("Successfully sent notice to telegram")
	return nil
}

// SendNoticeToEmail sends a notification to Email
func (sv *NoticeService) SendNoticeToEmail(ctx context.Context, email, message string) error {
	lg := sv.lg.With().Str("method", "SendNoticeToEmail").Logger()

	fn := func() error {
		err := sv.em.SendEmail([]string{email}, "Event Booker", message, false)
		if err != nil {
			lg.Warn().Err(err).Str("email", email).Msg("Failed to send notice to email")
		}
		return err
	}

	if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
		return fmt.Errorf("send notice to email after all attempts: %w", err)
	}

	lg.Debug().Str("email", email).Msg("Successfully sent notice to email")
	return nil
}
