package handler

import "github.com/gin-gonic/gin"

func (h *Handler) createEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "createEvent").Logger()

	log.Debug().Msg("handler is starting")
}
