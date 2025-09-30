package app

import (
	"fmt"

	"github.com/wb-go/wbf/zlog"
)

type App struct {
	cfg  *appConfig
	deps *dependencies
}

func New() (*App, error) {
	cfg, err := newAppConfig("./config/config.yaml", "./.env", "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error creating configuration")
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

	err = zlog.SetLevel(cfg.lg.Level)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "logger").Msg("error set level")
		return nil, fmt.Errorf("error set log level: %w", err)
	}

	return &App{
		cfg: cfg,
	}, nil
}

func (a *App) Run() error {
	// ------------------------- TEMP -------------------------
	fmt.Printf("port: %d\n", a.cfg.tr.TrHTTP.Srv.Port)
	fmt.Printf("logLevel: %s\n", a.cfg.lg.Level)
	return nil
}
