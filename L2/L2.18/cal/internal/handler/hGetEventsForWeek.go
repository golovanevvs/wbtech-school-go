package handler

import "github.com/gin-gonic/gin"

func (h *Handler) getEventsForWeek(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForWeek").Logger()

	log.Debug().Msg("handler is starting")
}
