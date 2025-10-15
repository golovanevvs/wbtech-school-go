package trhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	handler.IService
}

type HTTP struct {
	lg      zlog.Zerolog
	httpsrv *http.Server
}

func New(cfg *Config, sv IService) *HTTP {
	lg := zlog.Logger.With().Str("component", "transport-HTTP").Logger()
	return &HTTP{
		lg: lg,
		httpsrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: handler.New(cfg.Handler, sv).Rt,
		},
	}
}

func (h *HTTP) RunServer(cancel context.CancelFunc) {
	go func() {
		h.lg.Info().Str("addr", h.httpsrv.Addr).Msg("http server starting...")
		if err := h.httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.lg.Error().Err(err).Msg("error http server start")
			cancel()
		}
	}()
}

func (h *HTTP) ShutdownServer(ctx context.Context) error {
	h.lg.Info().Msg("http server stopping...")

	if err := h.httpsrv.Shutdown(ctx); err != nil {
		h.lg.Error().Err(err).Msg("error http server shutdown")
		return fmt.Errorf("failed http server shutdown: %w", err)
	}

	h.lg.Info().Msg("http server stopped successfully")

	return nil
}

func (h *HTTP) WaitForServer(host string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("%s/healthy", host))
		if err == nil && resp.StatusCode == http.StatusOK {
			h.lg.Info().Str("addr", h.httpsrv.Addr).Msg("http server started")
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	h.lg.Error().Dur("timeout", timeout).Msg("http server not ready")
	return fmt.Errorf("http server not ready after %s ms", timeout)
}
