package handler

import (
	"github.com/gin-gonic/gin"
	rediscache "github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/cache/redis"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/repository"
	"github.com/rs/zerolog"
)

type Handler struct {
	Router *gin.Engine
	logger *zerolog.Logger
	rp     *repository.Repository
	cache  *rediscache.RedisCache
}

func New(cfgHd *config.Handler, logger *zerolog.Logger, rp *repository.Repository, cache *rediscache.RedisCache) *Handler {
	log := logger.With().Str("component", "handler").Logger()
	gin.SetMode(cfgHd.GinMode)
	router := gin.New()
	return &Handler{
		Router: router,
		logger: &log,
		rp:     rp,
		cache:  cache,
	}
}

func (h *Handler) Run(addr string) error {
	return h.Router.Run(addr)
}

func (h *Handler) InitRoutes() {
	h.Router.Use(h.WithLogging())

	h.Router.GET("/order", h.hGetOrderByOrderUID)
}
