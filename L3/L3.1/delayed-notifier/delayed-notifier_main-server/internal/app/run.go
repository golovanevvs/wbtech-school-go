package app

import (
	"context"
	"fmt"
	"time"
)

func (a *App) Run(cancel context.CancelFunc) error {
	a.deps.tr.HTTP.RunServer(cancel)

	if err := a.deps.tr.HTTP.WaitForServer(fmt.Sprintf("%s/healthy", a.cfg.tr.TrHTTP.PublicHost), 30*time.Second); err != nil {
		return err
	}

	if err := a.deps.tg.SetWebhook(fmt.Sprintf("%s/telegram/webhook", a.cfg.tr.TrHTTP.PublicHost)); err != nil {
		return err
	}

	return nil
}
