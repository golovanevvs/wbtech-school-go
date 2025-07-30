package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	Router *gin.Engine
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Handler {
	return &Handler{
		Router: gin.New(),
		Logger: logger,
	}
}

func (h *Handler) Run(addr string) error {
	return h.Router.Run(addr)
}

func (h *Handler) InitRoutes() {
	h.Router.GET("/sample", h.hSample)
}
