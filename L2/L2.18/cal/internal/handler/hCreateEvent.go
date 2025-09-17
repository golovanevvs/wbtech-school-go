package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) createEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "createEvent").Logger()

	log.Debug().Msg("start handling create event request")

	if !strings.Contains(c.ContentType(), "application/json") {
		log.Error().Msg(errContentType)
		c.JSON(http.StatusBadRequest, model.Resp{Error: errContentType})
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

	if strings.TrimSpace(event.Title) == "" {
		log.Error().Msg(errEmptyTitle)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + errEmptyTitle,
		})
		return
	}

	if time.Time(event.Date).IsZero() {
		log.Error().Msg(errEmptyDate)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + errEmptyDate,
		})
		return
	}

	id := h.repository.Create(event)

	log.Debug().
		Str("user_id", event.UserId).
		Str("id", id).
		Str("title", event.Title).
		Str("comment", event.Comment).
		Time("date", time.Time(event.Date)).
		Msg("event created successfully")

	c.JSON(http.StatusOK, model.Resp{
		Id:     id,
		Result: "event created successfully",
	})
}
