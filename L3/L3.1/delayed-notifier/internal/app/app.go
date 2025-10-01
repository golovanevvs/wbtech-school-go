package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"
)

type App struct {
	cfg  *appConfig
	deps *dependencies
}

func New() (*App, error) {
	cfg, err := newConfig("./config/config.yaml", "./.env", "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error creating configuration")
		return nil, fmt.Errorf("error creating configuration: %w", err)
	}

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

func (a *App) Run() error {
	go a.deps.tr.HTTP.RunServer()
	return nil
}

func (a *App) GracefullShutdown(ctx context.Context, cancel context.CancelFunc) {
	lg := a.deps.lg.With().Str("component", "app").Logger()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		lg.Warn().Msg("Shutdown signal received, gracefully shutting down...")
		cancel()
	case <-ctx.Done():
	}

	a.deps.tr.HTTP.ShutdownServer(ctx)
}
