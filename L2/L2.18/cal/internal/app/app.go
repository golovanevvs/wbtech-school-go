package app

import (
	"fmt"
	"os"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/handler"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/server"
)

func Run() {
	t := 5

	zlog.Init()
	zlog.Logger.Info().Msg("starting cal...")

	envFile := ".env"
	envPath := fmt.Sprintf("./%s", envFile)

	cfg := config.New()
	err := cfg.Load(envPath)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to load configuration")
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}
	zlog.Logger.Info().Msg("configuration loaded successfully")

	logLevelStr := cfg.Logger.LogLevel
	logLevel, err := zlog.ParseLogLevel(logLevelStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse log level")
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}
	zlog.Logger.Info().Str("logLevel", logLevel.String()).Msg("logging level")
	zlog.Logger = zlog.Logger.Level(logLevel)

	hd := handler.New(&cfg.Handler, &zlog.Logger)
	hd.InitRoutes()

	srv := server.New(hd.Router, &zlog.Logger)

	srv.RunServerWithGracefulShutdown(cfg.Server.Addr, 5*time.Second)
}
