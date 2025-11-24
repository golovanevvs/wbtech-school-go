package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
)

// IBookingRp interface for booking repository
type IBookingRp interface {
	Create(booking *model.Booking) (*model.Booking, error)
	GetByID(id int) (*model.Booking, error)
	GetByUserID(userID int) ([]*model.Booking, error)
	GetByUserIDAndEventID(userID int, eventID int) (*model.Booking, error)
	GetByEventID(eventID int) ([]*model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id int) error
	GetExpiredBookings() ([]*model.Booking, error)
	UpdateStatus(id int, status model.BookingStatus) error
}

// EventRepo interface for event repository (needed for updating available places)
type IEventRpForBooking interface {
	GetByID(id int) (*model.Event, error)
	UpdateAvailablePlaces(eventID int, newAvailablePlaces int) error
}

// BookingService service for working with bookings
type BookingService struct {
	rp IBookingRp
	er IEventRpForBooking
}

// NewBookingService creates a new BookingService
func NewBookingService(rp IBookingRp, er IEventRpForBooking) *BookingService {
	return &BookingService{rp: rp, er: er}
}

// Create creates a new booking
func (sv *BookingService) Create(ctx context.Context, userID int, eventID int, bookingDeadlineMinutes int) (*model.Booking, error) {
	event, err := sv.er.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event.AvailablePlaces <= 0 {
		return nil, fmt.Errorf("no available places for event %d", eventID)
	}

	existingBooking, err := sv.rp.GetByUserIDAndEventID(userID, eventID)
	if err == nil {
		if existingBooking.Status == model.BookingCancelled {
			err = sv.rp.Delete(existingBooking.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete cancelled booking: %w", err)
			}
		} else {
			return nil, fmt.Errorf("booking already exists for user %d and event %d", userID, eventID)
		}
	}

	expiresAt := time.Now().Add(time.Duration(bookingDeadlineMinutes) * time.Minute)

	booking := &model.Booking{
		UserID:    userID,
		EventID:   eventID,
		Status:    model.BookingPending,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	createdBooking, err := sv.rp.Create(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	err = sv.er.UpdateAvailablePlaces(eventID, event.AvailablePlaces-1)
	if err != nil {
		_ = sv.rp.Delete(createdBooking.ID)
		return nil, fmt.Errorf("failed to update available places: %w", err)
	}

	return createdBooking, nil
}

// GetByID returns a booking by ID
func (sv *BookingService) GetByID(ctx context.Context, id int) (*model.Booking, error) {
	return sv.rp.GetByID(id)
}

// GetByUserID returns all bookings for a user
func (sv *BookingService) GetByUserID(ctx context.Context, userID int) ([]*model.Booking, error) {
	return sv.rp.GetByUserID(userID)
}

// GetByUserIDAndEventID returns a booking for a user and event
func (sv *BookingService) GetByUserIDAndEventID(ctx context.Context, userID int, eventID int) (*model.Booking, error) {
	return sv.rp.GetByUserIDAndEventID(userID, eventID)
}

// GetByEventID returns all bookings for an event
func (sv *BookingService) GetByEventID(ctx context.Context, eventID int) ([]*model.Booking, error) {
	return sv.rp.GetByEventID(eventID)
}

// Confirm confirms a booking
func (sv *BookingService) Confirm(ctx context.Context, bookingID int) error {
	booking, err := sv.rp.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking.Status != model.BookingPending {
		return fmt.Errorf("booking is not in pending status")
	}

	confirmedTime := time.Now()
	booking.Status = model.BookingConfirmed
	booking.ConfirmedAt = &confirmedTime

	err = sv.rp.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to confirm booking: %w", err)
	}

	return nil
}

// Cancel cancels a booking
func (sv *BookingService) Cancel(ctx context.Context, bookingID int) error {
	booking, err := sv.rp.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking.Status == model.BookingConfirmed {
		return fmt.Errorf("cannot cancel confirmed booking")
	}

	event, err := sv.er.GetByID(booking.EventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	cancelledTime := time.Now()
	booking.Status = model.BookingCancelled
	booking.CancelledAt = &cancelledTime

	err = sv.rp.Update(booking)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	err = sv.er.UpdateAvailablePlaces(booking.EventID, event.AvailablePlaces+1)
	if err != nil {
		return fmt.Errorf("failed to update available places: %w", err)
	}

	return nil
}

// ProcessExpiredBookings processes all expired bookings
func (sv *BookingService) ProcessExpiredBookings(ctx context.Context) error {
	expiredBookings, err := sv.rp.GetExpiredBookings()
	if err != nil {
		return fmt.Errorf("failed to get expired bookings: %w", err)
	}

	for _, booking := range expiredBookings {
		event, err := sv.er.GetByID(booking.EventID)
		if err != nil {
			continue
		}

		cancelledTime := time.Now()
		booking.Status = model.BookingCancelled
		booking.CancelledAt = &cancelledTime

		err = sv.rp.Update(booking)
		if err != nil {
			continue
		}

		err = sv.er.UpdateAvailablePlaces(booking.EventID, event.AvailablePlaces+1)
		if err != nil {
			continue
		}
	}

	return nil
}
