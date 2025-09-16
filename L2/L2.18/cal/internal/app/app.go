package app

import (
	"os"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/handler"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/server"
)

func Run() {
	zlog.Init()
	zlog.Logger.Info().Msg("starting cal...")

	cfg := config.New("./.env")

	err := zlog.SetLevel(cfg.Logger.LogLevel)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse log level")
		time.Sleep(cfg.App.DelayBeforeClosing)
		os.Exit(1)
	}

	hd := handler.New(cfg.Handler, &zlog.Logger)
	hd.InitRoutes()

	srv := server.New(cfg.Server, hd.Router, &zlog.Logger)

	srv.RunServerWithGracefulShutdown()
}
