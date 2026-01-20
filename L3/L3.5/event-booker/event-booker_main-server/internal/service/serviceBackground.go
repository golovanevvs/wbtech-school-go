package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgRetry"
	"github.com/wb-go/wbf/zlog"
)

// IBookingRpForBackground interface for booking repository in background service
type IBookingRpForBackground interface {
	GetExpiredBookings() ([]*model.Booking, error)
	UpdateStatus(id int, status model.BookingStatus) error
}

// IEventRpForBackground interface for event repository in background service
type IEventRpForBackground interface {
	GetByID(id int) (*model.Event, error)
	UpdateAvailablePlaces(eventID int, newAvailablePlaces int) error
}

// IUserRpForBackground interface for user repository in background service
type IUserRpForBackground interface {
	GetByID(id int) (*model.User, error)
}

// INoticeSvForBackground interface for notice service in background service
type INoticeSvForBackground interface {
	SendNotice(ctx context.Context, notice model.Notice)
}

// BackgroundService handles background tasks
type BackgroundService struct {
	lg        *zlog.Zerolog
	rs        *pkgRetry.Retry
	bookingRp IBookingRpForBackground
	eventRp   IEventRpForBackground
	userRp    IUserRpForBackground
	noticeSv  INoticeSvForBackground
}

// NewBackgroundService creates a new BackgroundService
func NewBackgroundService(parentLg *zlog.Zerolog, rs *pkgRetry.Retry, bookingRp IBookingRpForBackground, eventRp IEventRpForBackground, userRp IUserRpForBackground, noticeSv INoticeSvForBackground) *BackgroundService {
	lg := parentLg.With().Str("component", "service-BackgroundService").Logger()
	return &BackgroundService{
		lg:        &lg,
		rs:        rs,
		bookingRp: bookingRp,
		eventRp:   eventRp,
		userRp:    userRp,
		noticeSv:  noticeSv,
	}
}

// StartExpiredBookingProcessor starts the processor for expired bookings
func (sv *BackgroundService) StartExpiredBookingProcessor(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sv.processExpiredBookings(ctx); err != nil {
				sv.lg.Error().Err(err).Msg("Failed to process expired bookings")
			}
		}
	}
}

// processExpiredBookings processes all expired bookings
func (sv *BackgroundService) processExpiredBookings(ctx context.Context) error {
	expiredBookings, err := sv.bookingRp.GetExpiredBookings()
	if err != nil {
		return fmt.Errorf("failed to get expired bookings: %w", err)
	}

	for _, booking := range expiredBookings {
		event, err := sv.eventRp.GetByID(booking.EventID)
		if err != nil {
			sv.lg.Error().Err(err).Int("event_id", booking.EventID).Msg("Failed to get event for expired booking")
			continue
		}

		err = sv.bookingRp.UpdateStatus(booking.ID, model.BookingCancelled)
		if err != nil {
			sv.lg.Error().Err(err).Int("booking_id", booking.ID).Msg("Failed to update booking status")
			continue
		}

		err = sv.eventRp.UpdateAvailablePlaces(booking.EventID, event.AvailablePlaces+1)
		if err != nil {
			sv.lg.Error().Err(err).Int("event_id", booking.EventID).Msg("Failed to update available places")
			continue
		}

		sv.SendBookingExpiredNotice(ctx, booking, event)
	}

	return nil
}

// SendBookingExpiredNotice sends notification about expired booking
func (sv *BackgroundService) SendBookingExpiredNotice(ctx context.Context, booking *model.Booking, event *model.Event) {
	user, err := sv.userRp.GetByID(booking.UserID)
	if err != nil {
		return
	}

	message := fmt.Sprintf(
		"⏰ Срок брони истёк\n\n"+
			"Мероприятие: %s\n"+
			"Дата: %s\n\n"+
			"К сожалению, срок действия вашей брони истёк. Бронь была автоматически отменена.",
		event.Title,
		event.Date.Format("02.01.2006 в 15:04"),
	)

	notice := model.Notice{
		UserID:  user.ID,
		Message: message,
		Channels: model.NotificationChannels{
			Telegram: user.TelegramNotifications && user.TelegramChatID != nil,
			Email:    user.EmailNotifications,
		},
	}

	sv.noticeSv.SendNotice(ctx, notice)
}
