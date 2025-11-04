package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgPrometheus"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/getStatusHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	addNoticeHandler.IService
	deleteNoticeHandler.IService
	getStatusHandler.IService
	telegramHandler.IService
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

	addNoticeHandler := addNoticeHandler.New(&lg, rt, sv)
	addNoticeHandler.RegisterRoutes()

	deleteNoticeHandler := deleteNoticeHandler.New(&lg, rt, sv)
	deleteNoticeHandler.RegisterRoutes()

	getStatusHandler := getStatusHandler.New(&lg, rt, sv)
	getStatusHandler.RegisterRoutes()

	telegramHandler := telegramHandler.New(&lg, rt, sv)
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(&lg, rt)
	healthHandler.RegisterRoutes()

	return hd
}
