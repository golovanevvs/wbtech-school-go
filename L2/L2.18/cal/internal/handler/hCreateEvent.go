package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/customerrors"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) createEvent(c *gin.Context) {
	log := h.logger.With().Str("handler", "createEvent").Logger()

	log.Debug().Msg("start handling create event request")
	defer log.Debug().Msg("start handling create event request")

	if !strings.Contains(c.ContentType(), "application/json") {
		log.Error().Msg(customerrors.ErrContentTypeAJ.Error())
		c.JSON(http.StatusBadRequest, model.Resp{Error: customerrors.ErrContentTypeAJ.Error()})
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
		log.Error().Msg(customerrors.ErrEmptyUserID.Error())
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + customerrors.ErrEmptyUserID.Error(),
		})
		return
	}

	if strings.TrimSpace(event.Title) == "" {
		log.Error().Msg(customerrors.ErrEmptyTitle.Error())
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + customerrors.ErrEmptyTitle.Error(),
		})
		return
	}

	if time.Time(event.Date).IsZero() {
		log.Error().Msg(customerrors.ErrEmptyDate.Error())
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "validation failed: " + customerrors.ErrEmptyDate.Error(),
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
