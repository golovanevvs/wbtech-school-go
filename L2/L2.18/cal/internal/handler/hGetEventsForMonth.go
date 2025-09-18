package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) getEventsForMonth(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForMonth").Logger()

	log.Debug().Msg("start handling get events for month request")
	defer log.Debug().Msg("end handling get events for month request")

	userId := c.Query("user_id")
	monthStr := c.Query("month")
	if userId == "" || monthStr == "" {
		log.Error().Msg("missing user_id or month query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: missing user_id or month query parameter",
		})
		return
	}

	month, err := time.Parse("2006-01-02", monthStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse month query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: failed to parse month query parameter: " + err.Error(),
		})
		return
	}

	events, err := h.repository.LoadForMonth(userId, month)
	if err != nil {
		log.Error().Err(err).Msg("failed to get events for month")
		c.JSON(http.StatusInternalServerError, model.Resp{
			Error: "failed to get events for month: " + err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", userId).
		Str("month", monthStr).
		Int("events_count", len(events)).
		Msg("events for month retrieved successfully")

	c.JSON(http.StatusOK, model.Resp{
		Events: events,
		Result: "events for month retrieved successfully",
	})
}
