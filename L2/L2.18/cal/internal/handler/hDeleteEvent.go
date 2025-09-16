package handler

import "github.com/gin-gonic/gin"

func (h *Handler) deleteEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "deleteEvent").Logger()

	log.Debug().Msg("handler is starting")
}
