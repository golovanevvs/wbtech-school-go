package app

import (
	"context"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgErrors"
)

func (a *App) Run(cancel context.CancelFunc) error {
	// ctx := context.Background()

	a.deps.tr.HTTP.RunServer(cancel)
	time.Sleep(500 * time.Millisecond)

	if err := a.deps.tr.HTTP.WaitForServer(a.cfg.tr.TrHTTP.PublicHost); err != nil {
		return pkgErrors.Wrap(err, "failed to start http server")
	}

	return nil
}
