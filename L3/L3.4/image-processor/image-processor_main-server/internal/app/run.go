package app

import (
	"context"
	"fmt"
	"time"
)

func (a *App) Run(cancel context.CancelFunc) error {
	a.deps.tr.HTTP.RunServer(cancel)
	time.Sleep(500 * time.Millisecond)

	if err := a.deps.tr.HTTP.WaitForServer(a.cfg.tr.TrHTTP.PublicHost); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
