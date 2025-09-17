package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) createEvent(c *gin.Context) {
	if !strings.Contains(c.ContentType(), "application/json") {
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "content type must be application/json",
		})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "invalid json: " + err.Error(),
		})
		return
	}

	if strings.TrimSpace(event.UserId) == "" {
		c.JSON(http.StatusBadRequest, model.Resp{
			Error: "user_id must be not empty",
		})
		return
	}

	id := h.repository.Create(event)

	c.JSON(http.StatusOK, model.Resp{
		Id:     id,
		Result: "event created successfully",
	})
}
