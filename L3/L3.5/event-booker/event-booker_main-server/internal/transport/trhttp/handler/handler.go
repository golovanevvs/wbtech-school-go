package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgPrometheus"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/bookingHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/middleware"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler/telegramHandler"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	BookingService() bookingHandler.ISvForBookingHandler
	TelegramService() telegramHandler.ISvForTelegramHandler
	MiddlewareService() middleware.IServiceForAuthHandler
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

	telegramHandler := telegramHandler.New(&lg, rt, sv.TelegramService())
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(&lg, rt)
	healthHandler.RegisterRoutes()

	protected := rt.Group("/")
	protected.Use(authMiddleware.JWTMiddleware)
	{
		bookingHandler := bookingHandler.New(&lg, sv.BookingService())
		bookingHandler.RegisterRoutes(protected)
	}

	return hd
}
