package handler

import "github.com/gin-gonic/gin"

func (h *Handler) getEventsForDay(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForDay").Logger()

	log.Debug().Msg("handler is starting")
}
