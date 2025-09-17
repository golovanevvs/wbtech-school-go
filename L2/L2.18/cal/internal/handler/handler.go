package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/repository"
	"github.com/rs/zerolog"
)

const (
	errContentType = "content type must be application/json"
	errEmptyUserID = "user_id must not be empty"
	errEmptyID     = "id must not be empty"
	errEmptyTitle  = "title must not be empty"
	errEmptyDate   = "date must not be empty"
)

type Handler struct {
	Router     *gin.Engine
	logger     *zerolog.Logger
	repository *repository.Repository
}

func New(config *config.Handler, logger *zerolog.Logger, repository *repository.Repository) *Handler {
	log := logger.With().Str("component", "handler").Logger()
	gin.SetMode(config.GinMode)
	router := gin.New()

	return &Handler{
		Router:     router,
		logger:     &log,
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
