package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
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

// BackgroundService handles background tasks
type BackgroundService struct {
	lg        *zlog.Zerolog
	bookingRp IBookingRpForBackground
	eventRp   IEventRpForBackground
}

// NewBackgroundService creates a new BackgroundService
func NewBackgroundService(parentLg *zlog.Zerolog, bookingRp IBookingRpForBackground, eventRp IEventRpForBackground) *BackgroundService {
	lg := parentLg.With().Str("component", "service-BackgroundService").Logger()
	return &BackgroundService{
		lg:        &lg,
		bookingRp: bookingRp,
		eventRp:   eventRp,
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
	}

	return nil
}
