package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) getEventsForDay(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForDay").Logger()

	log.Debug().Msg("start handling get events for day request")
	defer log.Debug().Msg("start handling get events for day request")

	userId := c.Query("user_id")
	dayStr := c.Query("day")
	if userId == "" || dayStr == "" {
		log.Error().Msg("missing user_id or day query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: missing user_id or day query parameter",
		})
		return
	}

	day, err := time.Parse("2006-01-02", dayStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse day query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: failed to parse day query parameter: " + err.Error(),
		})
		return
	}

	events, err := h.repository.LoadForDay(userId, day)
	if err != nil {
		log.Error().Err(err).Msg("failed to get events for day")
		c.JSON(http.StatusInternalServerError, model.Resp{
			Error: "failed to get events for day: " + err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", userId).
		Str("day", dayStr).
		Int("events_count", len(events)).
		Msg("events for day retrieved successfully")

	c.JSON(http.StatusOK, model.Resp{
		Events: events,
		Result: "events for day retrieved successfully"})
}
