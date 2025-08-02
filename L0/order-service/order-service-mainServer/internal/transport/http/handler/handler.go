package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/config"
	"github.com/rs/zerolog"
)

type Handler struct {
	Router *gin.Engine
	Logger *zerolog.Logger
}

func New(cfgHd *config.Handler, logger *zerolog.Logger) *Handler {
	gin.SetMode(cfgHd.GinMode)
	router := gin.New()
	return &Handler{
		Router: router,
		Logger: logger,
	}
}

func (h *Handler) Run(addr string) error {
	return h.Router.Run(addr)
}

func (h *Handler) InitRoutes() {
	h.Router.Use(h.WithLogging())

	h.Router.GET("/sample", h.hSample)
}
