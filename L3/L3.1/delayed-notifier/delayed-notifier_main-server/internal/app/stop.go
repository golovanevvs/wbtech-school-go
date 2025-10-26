package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/zlog"
)

func (a *App) GracefullShutdown(ctx context.Context, cancel context.CancelFunc) {
	lg := zlog.Logger.With().Str("component", "app").Logger()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		lg.Warn().Msgf("%s shutdown signal received, gracefully shutting down...", pkgConst.Warn)
		cancel()
	case <-ctx.Done():
	}

	if err := a.deps.tr.HTTP.ShutdownServer(ctx); err != nil {
		lg.Error().Err(err).Msgf("%s failed to shutdown http server", pkgConst.Error)
	}

	if err := a.rm.closeAll(); err != nil {
		lg.Error().Err(err).Msgf("%s failed to close resources", pkgConst.Error)
	}

}
