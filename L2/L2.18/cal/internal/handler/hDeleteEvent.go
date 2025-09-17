package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) deleteEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "deleteEvent").Logger()

	log.Debug().Msg("start handling delete event request")

	if !strings.Contains(c.ContentType(), "application/json") {
		log.Error().Msg(errContentType)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: errContentType,
		})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid json: " + err.Error(),
		})
		return
	}

	if strings.TrimSpace(event.UserId) == "" {
		log.Error().Msg(errEmptyUserID)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + errEmptyUserID,
		})
		return
	}

	if strings.TrimSpace(event.Id) == "" {
		log.Error().Msg(errEmptyID)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + errEmptyID,
		})
		return
	}

	err := h.repository.Delete(event)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete event")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", event.UserId).
		Str("id", event.Id).
		Msg("event deleted successfully")

	c.JSON(http.StatusOK, model.Resp{
		Id:     event.Id,
		Result: "event deleted successfully",
	})
}
