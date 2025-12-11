package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPrometheus"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport/trhttp/handler/addShortURL"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport/trhttp/handler/getAnalytics"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport/trhttp/handler/getOriginalURL"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport/trhttp/handler/healthHandler"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	addShortURL.IAddShortURLService
	getOriginalURL.IGetOriginalURLService
	getOriginalURL.IAddClickEventService
	getAnalytics.IGetAnalyticsService
}

type Handler struct {
	Rt *ginext.Engine
}

func New(cfg *Config, parentLg *zlog.Zerolog, sv IService, publicHost string, webPublicHost string) *Handler {
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

	addShortURLHandler := addShortURL.New(&lg, rt, sv)
	addShortURLHandler.RegisterRoutes()

	getOriginalURLHandler := getOriginalURL.New(&lg, rt, sv, sv)
	getOriginalURLHandler.RegisterRoutes()

	getAnalyticsHandler := getAnalytics.New(&lg, rt, sv)
	getAnalyticsHandler.RegisterRoutes()

	healthHandler := healthHandler.New(&lg, rt)
	healthHandler.RegisterRoutes()

	return hd
}
