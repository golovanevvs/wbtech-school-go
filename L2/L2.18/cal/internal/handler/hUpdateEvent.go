package handler

import "github.com/gin-gonic/gin"

func (h *Handler) updateEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "updateEvent").Logger()

	log.Debug().Msg("handler is starting")
}
