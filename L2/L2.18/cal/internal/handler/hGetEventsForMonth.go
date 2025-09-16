package handler

import "github.com/gin-gonic/gin"

func (h *Handler) getEventsForMonth(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForMonth").Logger()

	log.Debug().Msg("handler is starting")
}
