package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgPrometheus"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/authHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/bookingHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/eventHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/middleware"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/telegramHandler"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	BookingService() bookingHandler.ISvForBookingHandler
	EventService() eventHandler.ISvForEventHandler
	TelegramService() telegramHandler.ISvForTelegramHandler
	MiddlewareService() middleware.ISvForAuthHandler
	AuthService() authHandler.ISvForAuthHandler
}

type Handler struct {
	Rt *ginext.Engine
}

func New(
	cfg *Config,
	parentLg *zlog.Zerolog,
	sv IService,
	publicHost string,
	webPublicHost string,
) *Handler {
	lg := parentLg.With().Str("component", "handler").Logger()

	rt := ginext.New(cfg.GinMode)

	rt.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			publicHost,
			webPublicHost,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	pkgPrometheus.Init()
	rt.Use(pkgPrometheus.GinMiddleware())

	hd := &Handler{
		Rt: rt,
	}

	authMiddleware := middleware.NewAuthMiddleware(parentLg, sv.MiddlewareService())

	// Публичные маршруты (без аутентификации)
	public := rt.Group("/")
	{
		eventHandler := eventHandler.New(&lg, sv.EventService())
		eventHandler.RegisterPublicRoutes(public) // GET /events, GET /events/:id
	}

	// Защищенные маршруты (требуют аутентификации)
	protected := rt.Group("/")
	protected.Use(authMiddleware.JWTMiddleware)
	{
		// Защищенные маршруты мероприятий
		eventHandler := eventHandler.New(&lg, sv.EventService())
		eventHandler.RegisterProtectedRoutes(protected) // POST /events, PUT /events/:id, DELETE /events/:id

		// Маршруты бронирования
		bookingHandler := bookingHandler.New(&lg, sv.BookingService())
		bookingHandler.RegisterRoutes(protected) // POST /bookings, POST /bookings/:id/confirm и т.д.

		// Защищенные маршруты аутентификации (только для получения текущего пользователя и обновления)
		authHandler := authHandler.New(&lg, sv.AuthService())
		authHandler.RegisterProtectedRoutes(protected) // GET /auth/me, PUT /auth/update
	}

	// Публичные маршруты аутентификации (регистрация, вход, обновление токена)
	publicAuth := rt.Group("/auth")
	{
		authHandler := authHandler.New(&lg, sv.AuthService())
		authHandler.RegisterPublicRoutes(publicAuth) // POST /auth/register, POST /auth/login, POST /auth/refresh
	}

	// Обработчики, которые не требуют аутентификации
	telegramHandler := telegramHandler.New(&lg, rt, sv.TelegramService())
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(&lg, rt)
	healthHandler.RegisterRoutes()

	return hd
}
