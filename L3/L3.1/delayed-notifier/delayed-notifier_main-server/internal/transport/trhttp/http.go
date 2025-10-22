package trhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	handler.IService
}

type HTTP struct {
	lg                         *zlog.Zerolog
	httpsrv                    *http.Server
	retryStrategyForWaitServer retry.Strategy
}

func New(cfg *Config, parentLg *zlog.Zerolog, sv IService) *HTTP {
	lg := parentLg.With().Str("component", "HTTP").Logger()
	return &HTTP{
		lg: &lg,
		httpsrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: handler.New(cfg.Handler, &lg, sv).Rt,
		},
		retryStrategyForWaitServer: retry.Strategy(cfg.RetryStrategyForWaitServer),
	}
}

func (h *HTTP) RunServer(cancel context.CancelFunc) {
	go func() {
		h.lg.Debug().Str("addr", h.httpsrv.Addr).Msgf("%s http server starting...", pkgConst.Start)
		if err := h.httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.lg.Error().Err(err).Str("addr", h.httpsrv.Addr).Msgf("%s error http server start", pkgConst.Error)
			cancel()
		}
	}()
}

func (h *HTTP) ShutdownServer(ctx context.Context) error {
	h.lg.Debug().Str("addr", h.httpsrv.Addr).Msgf("%s http server stopping...", pkgConst.Start)

	if err := h.httpsrv.Shutdown(ctx); err != nil {
		pkgErrors.Wrapf(err, "http server shutdown, address: %s", h.httpsrv.Addr)
	}

	h.lg.Info().Str("addr", h.httpsrv.Addr).Msgf("%s http server stopped successfully", pkgConst.Info)

	return nil
}

func (h *HTTP) WaitForServer(host string) error {
	fn := func() error {
		resp, err := http.Get(fmt.Sprintf("%s/healthy", host))
		if err != nil {
			h.lg.Warn().Err(err).Str("addr", h.httpsrv.Addr).Msgf("%s failed to start http server", pkgConst.Warn)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			h.lg.Warn().Err(err).Str("addr", h.httpsrv.Addr).Msgf("%s failed to start http server", color.YellowString("⚠"))
			return err
		}
		return nil
	}

	if err := retry.Do(fn, h.retryStrategyForWaitServer); err != nil {
		return pkgErrors.Wrapf(err, "start http server, address: %s, attempts: %d", h.httpsrv.Addr, h.retryStrategyForWaitServer.Attempts)
	}

	h.lg.Info().Str("addr", h.httpsrv.Addr).Msg("ℹ http server started successfully")

	return nil
}
