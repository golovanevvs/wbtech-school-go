package handler

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/wb-go/wbf/ginext"
)

type IService interface {
	addNoticeHandler.IService
	deleteNoticeHandler.IService
}

type Handler struct {
	Rt *ginext.Engine
}

func New(cfg *Config, sv IService) *Handler {
	rt := ginext.New(cfg.GinMode)
	hd := &Handler{
		Rt: rt,
	}

	addNoticeHandler := addNoticeHandler.New(rt, sv)
	addNoticeHandler.RegisterRoutes()
	deleteNoticeHandler := deleteNoticeHandler.New(rt, sv)
	deleteNoticeHandler.RegisterRoutes()

	return hd
}
