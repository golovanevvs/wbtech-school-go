package handler

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
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

func New(cfg *Config, parentLg *zlog.Zerolog, sv IService, publicHost string) *Handler {
	lg := parentLg.With().Str("component", "handler").Logger()

	rt := ginext.New(cfg.GinMode)

	rt.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			fmt.Sprintf("http://localhost:%d", cfg.WebClientPort),
			fmt.Sprintf("http://127.0.0.1:%d", cfg.WebClientPort),
			publicHost,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	hd := &Handler{
		Rt: rt,
	}

	addNoticeHandler := addNoticeHandler.New(&lg, rt, sv)
	addNoticeHandler.RegisterRoutes()

	deleteNoticeHandler := deleteNoticeHandler.New(rt, sv)
	deleteNoticeHandler.RegisterRoutes()

	getStatusHandler := getStatusHandler.New(&lg, rt, sv)
	getStatusHandler.RegisterRoutes()

	telegramHandler := telegramHandler.New(rt, sv)
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(rt)
	healthHandler.RegisterRoutes()

	return hd
}
