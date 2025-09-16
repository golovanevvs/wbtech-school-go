package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/rs/zerolog"
)

type Handler struct {
	Router *gin.Engine
	logger *zerolog.Logger
}

func New(config *config.Handler, logger *zerolog.Logger) *Handler {
	log := logger.With().Str("component", "handler").Logger()
	gin.SetMode(config.GinMode)
	router := gin.New()

	return &Handler{
		Router: router,
		logger: &log,
	}
}

func (h *Handler) InitRoutes() {
	h.Router.Use(h.WithLogging())

	h.Router.GET("/sample", h.hSample)

}
