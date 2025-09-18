package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/repository"
	"github.com/rs/zerolog"
)

type Handler struct {
	Router     *gin.Engine
	logger     zerolog.Logger
	repository *repository.Repository
}

func New(config *config.Handler, logger *zerolog.Logger, repository *repository.Repository) *Handler {
	gin.SetMode(config.GinMode)
	router := gin.New()

	return &Handler{
		Router:     router,
		logger:     logger.With().Str("component", "handler").Logger(),
		repository: repository,
	}
}

func (h *Handler) InitRoutes() {
	h.Router.Use(h.WithLogging())

	h.Router.POST("/event/create", h.createEvent)
	h.Router.POST("/event/update", h.updateEvent)
	h.Router.POST("/event/delete", h.deleteEvent)
	h.Router.GET("/events/day", h.getEventsForDay)
	h.Router.GET("/events/week", h.getEventsForWeek)
	h.Router.GET("/events/month", h.getEventsForMonth)
}
