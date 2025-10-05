package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"
)

func (a *App) GracefullShutdown(ctx context.Context, cancel context.CancelFunc) {
	lg := zlog.Logger.With().Str("component", "app").Logger()
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
