package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/repository"
)

// Service structure that combines all services
type Service struct {
	User       *UserService
	Event      *EventService
	Booking    *BookingService
	Auth       *AuthService
	Background *BackgroundService
}

// New creates a new Service structure
func New(cfg *Config, rp *repository.Repository) *Service {
	userService := NewUserService(rp.User())
	eventService := NewEventService(rp.Event())
	bookingService := NewBookingService(rp.Booking())
	authService := NewAuthService(rp, cfg)
	backgroundService := NewBackgroundService(rp)

	return &Service{
		User:       userService,
		Event:      eventService,
		Booking:    bookingService,
		Auth:       authService,
		Background: backgroundService,
	}
}
