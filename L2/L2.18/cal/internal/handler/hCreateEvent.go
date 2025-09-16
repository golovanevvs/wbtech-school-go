package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

func (h *Handler) createEvent(c *gin.Context) {
	if !strings.Contains(c.ContentType(), "application/json") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "content type must be application/json",
		})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json: " + err.Error(),
		})
		return
	}

	h.repository.Create(event)

	c.JSON(http.StatusOK, gin.H{
		"result": "event created successfully",
	})
}
