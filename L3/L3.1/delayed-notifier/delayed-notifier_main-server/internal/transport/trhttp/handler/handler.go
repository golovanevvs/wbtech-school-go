package handler

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	addNoticeHandler.IService
	deleteNoticeHandler.IService
	telegramHandler.IService
}

type Handler struct {
	Rt *ginext.Engine
}

func New(cfg *Config, parentLg *zlog.Zerolog, sv IService) *Handler {
	lg := parentLg.With().Str("component-2", "handler").Logger()

	rt := ginext.New(cfg.GinMode)
	hd := &Handler{
		Rt: rt,
	}

	addNoticeHandler := addNoticeHandler.New(&lg, rt, sv)
	addNoticeHandler.RegisterRoutes()

	deleteNoticeHandler := deleteNoticeHandler.New(rt, sv)
	deleteNoticeHandler.RegisterRoutes()

	telegramHandler := telegramHandler.New(rt, sv)
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(rt)
	healthHandler.RegisterRoutes()

	return hd
}
