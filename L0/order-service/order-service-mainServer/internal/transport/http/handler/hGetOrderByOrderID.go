package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	rediscache "github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/cache/redis"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/repository/postgres"
)

func (h *Handler) hGetOrderByOrderUID(c *gin.Context) {
	log := h.logger.With().Str("handler", "hGetOrderByOrderID").Logger()

	log.Debug().Msg("handler is starting")

	orderUID := c.Request.FormValue("order_uid")
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order UID is required"})
		return
	}

	ctx := c.Request.Context()

	c.Header("Content-Type", "application/json; charset=UTF-8")

	order, err := h.cache.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, rediscache.ErrRedisNil) {
			order, err = h.rp.GetOrderByOrderUID(ctx, orderUID)
			if err != nil {
				if errors.Is(err, postgres.ErrOrderNotFound) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "order not found in database"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				c.JSON(http.StatusOK, order)
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, order)
}
