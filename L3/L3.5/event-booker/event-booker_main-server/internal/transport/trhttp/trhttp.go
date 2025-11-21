package trhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	handler.IService
}

type HTTP struct {
	lg      *zlog.Zerolog
	rs      *pkgRetry.Retry
	httpsrv *http.Server
}

func New(
	cfg *Config,
	parentLg *zlog.Zerolog,
	rs *pkgRetry.Retry,
	sv IService,
) *HTTP {
	lg := parentLg.With().Str("component", "HTTP").Logger()
	return &HTTP{
		lg: &lg,
		rs: rs,
		httpsrv: &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.Port),
			Handler: handler.New(
				cfg.Handler,
				&lg,
				sv,
				cfg.PublicHost,
				cfg.WebPublicHost,
			).Rt,
		},
	}
}

func (h *HTTP) RunServer(cancel context.CancelFunc) {
	go func() {
		h.lg.Info().Str("addr", h.httpsrv.Addr).Msgf("%s http server starting...", pkgConst.Starting)
		if err := h.httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.lg.Error().Err(err).Str("addr", h.httpsrv.Addr).Msgf("%s error http server start", pkgConst.Error)
			cancel()
		}
	}()
}

func (h *HTTP) ShutdownServer(ctx context.Context) error {
	h.lg.Debug().Str("addr", h.httpsrv.Addr).Msgf("%s http server stopping...", pkgConst.Starting)

	if err := h.httpsrv.Shutdown(ctx); err != nil {
		pkgErrors.Wrapf(err, "http server shutdown, address: %s", h.httpsrv.Addr)
	}

	h.lg.Info().Str("addr", h.httpsrv.Addr).Msgf("%s http server stopped successfully", pkgConst.Finished)

	return nil
}

func (h *HTTP) WaitForServer(host string) error {
	fn := func() error {
		resp, err := http.Get(fmt.Sprintf("%s/healthy", host))
		if err != nil {
			h.lg.Warn().Err(err).Str("addr", h.httpsrv.Addr).Str("host", host).Msgf("%s failed to start http server", pkgConst.Warn)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			h.lg.Warn().Err(err).Str("addr", h.httpsrv.Addr).Msgf("%s failed to start http server", pkgConst.Warn)
			return err
		}
		return nil
	}

	if err := retry.Do(fn, retry.Strategy(*h.rs)); err != nil {
		return pkgErrors.Wrapf(err, "start http server, address: %s, attempts: %d", h.httpsrv.Addr, h.rs.Attempts)
	}

	h.lg.Info().Str("addr", h.httpsrv.Addr).Msgf("%s http server started successfully", pkgConst.Finished)

	return nil
}
