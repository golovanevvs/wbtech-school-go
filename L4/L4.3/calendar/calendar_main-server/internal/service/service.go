package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/authHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/bookingHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/eventHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/middleware"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/wb-go/wbf/zlog"
)

// Service structure that combines all services
type Service struct {
	User          *UserService
	Event         *EventService
	Booking       *BookingService
	Auth          *AuthService
	Background    *BackgroundService
	TelegramStart *TelegramStartService
	Notice        *NoticeService
}

// New creates a new Service structure
func New(
	cfg *Config,
	rp *repository.Repository,
	tgClient *pkgTelegram.Client,
	emailClient *pkgEmail.Client,
	rs *pkgRetry.Retry,
) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()

	userService := NewUserService(rp.User())
	eventService := NewEventService(rp.Event())
	bookingService := NewBookingService(rp.Booking(), rp.Event())
	authService := NewAuthService(cfg, rp.User(), rp.RefreshToken())
	backgroundService := NewBackgroundService(&lg, rp.Booking(), rp.Event())
	telegramStartService := NewTelegramStartService(tgClient, rp.User())
	noticeService := NewNoticeService(&lg, rs, tgClient, emailClient, rp.NoticeRepository)

	return &Service{
		User:          userService,
		Event:         eventService,
		Booking:       bookingService,
		Auth:          authService,
		Background:    backgroundService,
		TelegramStart: telegramStartService,
		Notice:        noticeService,
	}
}

// BookingService returns the booking service
func (sv *Service) BookingService() bookingHandler.ISvForBookingHandler {
	return sv.Booking
}

// TelegramService returns the telegram service
func (sv *Service) TelegramService() telegramHandler.ISvForTelegramHandler {
	return sv.TelegramStart
}

// MiddlewareService returns the auth service for middleware
func (sv *Service) MiddlewareService() middleware.ISvForAuthHandler {
	return sv.Auth
}

// AuthService returns the auth service
func (sv *Service) AuthService() authHandler.ISvForAuthHandler {
	return sv.Auth
}

// EventService returns the event service
func (sv *Service) EventService() eventHandler.ISvForEventHandler {
	return sv.Event
}
