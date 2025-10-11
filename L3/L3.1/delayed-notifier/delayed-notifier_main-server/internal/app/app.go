package app

import (
	"fmt"

	"github.com/wb-go/wbf/zlog"
)

type App struct {
	cfg  *Config
	deps *dependencies
}

func New() (*App, error) {
	cfg, err := newConfig()
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error creating configuration")
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

	zlog.Logger.Info().Msg(cfg.String())

	zlog.Logger.Info().Str("component", "app").Msg("starting dependency initialization...")
	deps, err := newDependencyBuilder(cfg).build()
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error dependencies initialization")
		return nil, fmt.Errorf("error dependencies initialization: %w", err)
	}
	zlog.Logger.Info().Str("component", "app").Msg("dependencies have been initialized successfully")

	return &App{
		cfg:  cfg,
		deps: deps,
	}, nil
}
