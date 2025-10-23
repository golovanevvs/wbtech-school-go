package app

import (
	"fmt"

	"github.com/wb-go/wbf/zlog"
)

type App struct {
	cfg  *Config
	lg   *zlog.Zerolog
	deps *dependencies
	rm   *resourceManager
}

func New(env string) (*App, error) {
	lg := zlog.Logger.With().Str("component", "app").Logger()

	cfg, err := newConfig(env)
	if err != nil {
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

	lg.Info().Msg(cfg.String())

	lg.Info().Msg("starting dependency initialization...")
	deps, rm, err := newDependencyBuilder(cfg, &lg).build()
	if err != nil {
		return nil, fmt.Errorf("error dependencies initialization: %w", err)
	}
	lg.Info().Msg("dependencies have been initialized successfully")

	return &App{
		cfg:  cfg,
		lg:   &lg,
		deps: deps,
		rm:   rm,
	}, nil
}
