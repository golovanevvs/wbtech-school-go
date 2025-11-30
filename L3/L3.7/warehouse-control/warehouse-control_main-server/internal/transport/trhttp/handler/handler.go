package handler

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPrometheus"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/authHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/healthHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/itemHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/middleware"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	MiddlewareService() middleware.ISvForAuthHandler
	AuthService() authHandler.ISvForAuthHandler
	ItemService() service.IItemService
}

type Handler struct {
	Rt *ginext.Engine
}

func New(
	cfg *Config,
	parentLg *zlog.Zerolog,
	sv IService,
	publicHost string,
	webPublicHost string,
	accessTokenExp, refreshTokenExp time.Duration,
) *Handler {
	lg := parentLg.With().Str("component", "handler").Logger()

	rt := ginext.New(cfg.GinMode)

	rt.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			publicHost,
			webPublicHost,
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	pkgPrometheus.Init()
	rt.Use(pkgPrometheus.GinMiddleware())

	hd := &Handler{
		Rt: rt,
	}

	authMiddleware := middleware.NewAuthMiddleware(parentLg, sv.MiddlewareService(), accessTokenExp, refreshTokenExp)

	// public := rt.Group("/")
	{

	}

	protected := rt.Group("/")
	protected.Use(authMiddleware.JWTMiddleware)
	{
		authHandler := authHandler.New(&lg, sv.AuthService(), publicHost, accessTokenExp, refreshTokenExp)
		authHandler.RegisterProtectedRoutes(protected)

		// Добавляем хендлер для товаров
		itemHandler := itemHandler.New(&lg, sv.ItemService())
		itemHandler.RegisterProtectedRoutes(protected)
	}

	publicAuth := rt.Group("/auth")
	{
		authHandler := authHandler.New(&lg, sv.AuthService(), publicHost, accessTokenExp, refreshTokenExp)
		authHandler.RegisterPublicRoutes(publicAuth)
	}

	healthHandler := healthHandler.New(&lg, rt)
	healthHandler.RegisterRoutes()

	return hd
}
