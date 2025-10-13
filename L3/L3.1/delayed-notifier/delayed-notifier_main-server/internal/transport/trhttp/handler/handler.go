package handler

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/wb-go/wbf/ginext"
)

type IService interface {
	AddNotice() addNoticeHandler.IService
	DeleteNotice() deleteNoticeHandler.IService
	TelegramHandler() telegramHandler.IService
}

type Handler struct {
	Rt *ginext.Engine
}

func New(cfg *Config, sv IService) *Handler {
	rt := ginext.New(cfg.GinMode)
	hd := &Handler{
		Rt: rt,
	}

	addNoticeHandler := addNoticeHandler.New(rt, sv.AddNotice())
	addNoticeHandler.RegisterRoutes()

	deleteNoticeHandler := deleteNoticeHandler.New(rt, sv.DeleteNotice())
	deleteNoticeHandler.RegisterRoutes()

	telegramHandler := telegramHandler.New(rt, sv.TelegramHandler())
	telegramHandler.RegisterRoutes()

	healthHandler := healthHandler.New(rt)
	healthHandler.RegisterRoutes()

	return hd
}
