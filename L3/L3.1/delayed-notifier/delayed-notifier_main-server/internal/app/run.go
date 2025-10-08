package app

import (
	"context"
	"fmt"
)

func (a *App) Run(cancel context.CancelFunc) error {
	fmt.Println("запуск")
	a.deps.tr.HTTP.RunServer(cancel)
	return nil
}
