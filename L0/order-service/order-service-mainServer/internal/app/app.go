package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/handler"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/logger/zlog"
	"github.com/rs/zerolog"
)

type app struct {
	config  *config.Config
	handler *handler.Handler
	logger  *zerolog.Logger
}

func Run() {
	zlog.Init()

	configFile := "config.yaml"
	envFile := ".env"
	configPath := fmt.Sprintf("../../config/%s", configFile)
	configDefaultFile := "default-config.yaml"
	configDefaultPath := fmt.Sprintf("../../config/%s", configDefaultFile)
	envPath := fmt.Sprintf("../../%s", envFile)

	zlog.Logger.Info().Msg("Starting order-service-mainServer...")
	zlog.Logger.Info().Str("file", configFile).Msg("Loading configuration...")

	cfg := config.New()
	err := cfg.Load(configPath, envPath, "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("file", configFile).Msg("Failed to load configuration")
		zlog.Logger.Warn().Str("file", configDefaultFile).Msg("Loading default configuration...")
		err := cfg.Load(configDefaultPath, envPath, "")
		if err != nil {
			zlog.Logger.Fatal().Err(err).Msg("Failed to load default configuration")
		}
	}
	zlog.Logger.Info().Msg("Configuration loaded successfully")

	logLevelStr := cfg.Logger.LogLevel
	logLevel, err := zlog.ParseLogLevel(logLevelStr)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to parse log level")
	}
	zlog.Logger = zlog.Logger.Level(logLevel)
	zlog.Logger.Info().Str("logLevel", zlog.Logger.GetLevel().String()).Msg("Logging level")

	h := handler.New(&zlog.Logger)

	a := app{
		config:  cfg,
		handler: h,
		logger:  &zlog.Logger,
	}

	a.handler.InitRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info().Str("address", a.config.Server.Addr).Msg("Starting HTTP server")
		if err := a.handler.Run(a.config.Server.Addr); err != nil {
			a.logger.Error().Err(err).Msg("Failed to start HTTP server")
			cancel()
		}
	}()

	select {
	case sig := <-quit:
		a.logger.Warn().Str("signal", sig.String()).Msg("Received shutdown signal")
	case <-ctx.Done():
		a.logger.Warn().Msg("Context cancelled")
	}

	a.logger.Warn().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv := &http.Server{
		Addr:    a.config.Server.Addr,
		Handler: a.handler.Router,
	}

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	a.logger.Info().Msg("Server exited")
}
