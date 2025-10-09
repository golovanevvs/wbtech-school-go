package app

import (
	"context"
)

func (a *App) Run(cancel context.CancelFunc) error {
	a.deps.tr.HTTP.RunServer(cancel)
	return nil
}
