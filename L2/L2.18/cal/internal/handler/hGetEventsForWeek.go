package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) getEventsForWeek(c *gin.Context) {
	log := h.logger.With().Str("handler", "getEventsForWeek").Logger()

	log.Debug().Msg("start handling get events for week request")
	defer log.Debug().Msg("end handling get events for week request")

	userId := c.Query("user_id")
	weekStr := c.Query("week")
	if userId == "" || weekStr == "" {
		log.Error().Msg("missing user_id or week query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: missing user_id or week query parameter",
		})
		return
	}

	week, err := time.Parse("2006-01-02", weekStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse week query parameter")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid request: failed to parse week query parameter: " + err.Error(),
		})
		return
	}

	events, err := h.repository.LoadForWeek(userId, week)
	if err != nil {
		log.Error().Err(err).Msg("failed to get events for week")
		c.JSON(http.StatusInternalServerError, model.Resp{
			Error: "failed to get events for week: " + err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", userId).
		Str("week", weekStr).
		Int("events_count", len(events)).
		Msg("events for week retrieved successfully")

	c.JSON(http.StatusOK, model.Resp{
		Events: events,
		Result: "events for week retrieved successfully",
	})
}
