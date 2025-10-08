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
	cfg, err := newConfig()
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error creating configuration")
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

	zlog.Logger.Debug().Msg(cfg.String())

	deps, err := newDependencyBuilder(cfg).build()
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error dependencies initialization")
		return nil, fmt.Errorf("error dependencies initialization: %w", err)
	}

	return &App{
		cfg:  cfg,
		deps: deps,
	}, nil
}
