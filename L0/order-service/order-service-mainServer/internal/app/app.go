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
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/transport/http/handler"
	"github.com/rs/zerolog"
)

type app struct {
	config  *config.Config
	handler *handler.Handler
	logger  *zerolog.Logger
	rp      *repository.Repository
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zlog.Init()

	configFile := "config.yaml"
	envFile := ".env"
	configPath := fmt.Sprintf("./config/%s", configFile)
	configDefaultFile := "default-config.yaml"
	configDefaultPath := fmt.Sprintf("./config/%s", configDefaultFile)
	envPath := fmt.Sprintf("./%s", envFile)

	zlog.Logger.Info().Msg("starting order-service-mainServer...")
	zlog.Logger.Info().Str("file", configFile).Msg("loading configuration...")

	cfg := config.New()
	err := cfg.Load(configPath, envPath, "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("file", configFile).Msg("failed to load configuration")
		zlog.Logger.Warn().Str("file", configDefaultFile).Msg("loading default configuration...")
		err := cfg.Load(configDefaultPath, envPath, "")
		if err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to load default configuration")
		}
	}
	zlog.Logger.Info().Msg("configuration loaded successfully")

	logLevelStr := cfg.Logger.LogLevel
	logLevel, err := zlog.ParseLogLevel(logLevelStr)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to parse log level")
	}
	zlog.Logger.Info().Str("logLevel", logLevel.String()).Msg("logging level")
	zlog.Logger = zlog.Logger.Level(logLevel)

	rp, err := repository.New(ctx, &cfg.Repository, &zlog.Logger)
	if err != nil {
		os.Exit(1)
	}

	h := handler.New(&cfg.Handler, &zlog.Logger)

	a := app{
		config:  cfg,
		handler: h,
		logger:  &zlog.Logger,
		rp:      rp,
	}

	a.handler.InitRoutes()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info().Str("address", a.config.Server.Addr).Msg("starting HTTP server")
		if err := a.handler.Run(a.config.Server.Addr); err != nil {
			a.logger.Error().Err(err).Msg("failed to start HTTP server")
			cancel()
		}
	}()

	select {
	case sig := <-quit:
		a.logger.Warn().Str("signal", sig.String()).Msg("received shutdown signal")
	case <-ctx.Done():
		a.logger.Warn().Msg("context cancelled")
	}

	a.logger.Warn().Msg("shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv := &http.Server{
		Addr:    a.config.Server.Addr,
		Handler: a.handler.Router,
	}

	rp.Close()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	a.logger.Info().Msg("server exited")
	time.Sleep(5 * time.Second)
}
