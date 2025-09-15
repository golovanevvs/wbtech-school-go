package handler

import "github.com/gin-gonic/gin"

func (h *Handler) hSample(c *gin.Context) {
	log := h.logger.With().Str("handler", "hSample").Logger()

	log.Debug().Msg("handler is starting")
}
