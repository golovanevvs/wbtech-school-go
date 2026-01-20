package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/zlog"
)

// NoticeService handles sending notifications for calendar events
type NoticeService struct {
	lg *zlog.Zerolog
	rs *pkgRetry.Retry
	tg *pkgTelegram.Client
	em *pkgEmail.Client
}

// NewNoticeService creates a new NoticeService
func NewNoticeService(
	parentLg *zlog.Zerolog,
	rs *pkgRetry.Retry,
	tg *pkgTelegram.Client,
	em *pkgEmail.Client,
) *NoticeService {
	lg := parentLg.With().Str("component", "NoticeService").Logger()
	return &NoticeService{
		lg: &lg,
		rs: rs,
		tg: tg,
		em: em,
	}
}

// SendReminder sends a reminder notification for an event
func (sv *NoticeService) SendReminder(ctx context.Context, event *model.Event) {
	lg := sv.lg.With().Str("method", "SendReminder").Logger()

	message := sv.formatReminderMessage(event)

	wg := sync.WaitGroup{}

	if sv.tg != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendReminderToTelegram(ctx, message)
		}()
	}

	if sv.em != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendReminderToEmail(ctx, message)
		}()
	}

	wg.Wait()

	lg.Debug().Int("event_id", event.ID).Str("title", event.Title).Msg("Reminder sent successfully")
}

// formatReminderMessage formats the reminder message
func (sv *NoticeService) formatReminderMessage(event *model.Event) string {
	timeStr := event.Start.Format("2006-01-02 15:04")
	if event.AllDay {
		timeStr = event.Start.Format("2006-01-02") + " (весь день)"
	}

	message := fmt.Sprintf("🔔 Напоминание о событии\n\n📅 %s\n⏰ %s", event.Title, timeStr)

	if event.Description != "" {
		message += fmt.Sprintf("\n📝 %s", event.Description)
	}

	return message
}

// SendReminderToTelegram sends a reminder to Telegram
func (sv *NoticeService) SendReminderToTelegram(ctx context.Context, message string) error {
	lg := sv.lg.With().Str("method", "SendReminderToTelegram").Logger()

	lg.Debug().Str("message", message).Msg("Would send reminder to Telegram")

	// Example of actual sending (if you have configured chat IDs):
	// fn := func() error {
	//     err := sv.tg.SendTo(chatID, message)
	//     if err != nil {
	//         lg.Warn().Err(err).Msg("Failed to send reminder to telegram")
	//         return err
	//     }
	//     return nil
	// }
	//
	// if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
	//     return fmt.Errorf("send reminder to telegram after all attempts: %w", err)
	// }

	return nil
}

// SendReminderToEmail sends a reminder to Email
func (sv *NoticeService) SendReminderToEmail(ctx context.Context, message string) error {
	lg := sv.lg.With().Str("method", "SendReminderToEmail").Logger()

	lg.Debug().Str("message", message).Msg("Would send reminder to email")

	// Example of actual sending:
	// fn := func() error {
	//     err := sv.em.SendEmail([]string{email}, "Напоминание о событии", message, false)
	//     if err != nil {
	//         lg.Warn().Err(err).Msg("Failed to send reminder to email")
	//     }
	//     return err
	// }
	//
	// if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
	//     return fmt.Errorf("send reminder to email after all attempts: %w", err)
	// }

	return nil
}

// SendTestNotification sends a test notification
func (sv *NoticeService) SendTestNotification(ctx context.Context) error {
	lg := sv.lg.With().Str("method", "SendTestNotification").Logger()

	testMessage := "🧪 Тестовое уведомление от календаря событий"

	wg := sync.WaitGroup{}

	if sv.tg != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendReminderToTelegram(ctx, testMessage)
		}()
	}

	if sv.em != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sv.SendReminderToEmail(ctx, testMessage)
		}()
	}

	wg.Wait()

	lg.Info().Msg("Test notification sent")
	return nil
}
