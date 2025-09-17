package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) updateEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "updateEvent").Logger()

	log.Debug().Msg("start handling update event request")

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

	if strings.TrimSpace(event.Id) == "" {
		log.Error().Msg(errEmptyID)
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + errEmptyID,
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

	err := h.repository.Update(event)
	if err != nil {
		log.Error().Err(err).Msg("failed to update event")
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", event.UserId).
		Str("id", event.Id).
		Str("title", event.Title).
		Str("comment", event.Comment).
		Time("date", time.Time(event.Date)).
		Msg("event updated successfully")

	c.JSON(http.StatusOK, model.Resp{
		Id:     event.Id,
		Result: "event updated successfully",
	})
}
