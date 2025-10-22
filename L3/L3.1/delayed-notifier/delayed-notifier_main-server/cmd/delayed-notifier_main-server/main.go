package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/app"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	env := os.Getenv("ENV")
	if env == "local" {
		zlog.InitConsole()
	} else {
		zlog.Init()
	}

	lg := zlog.Logger.With().Str("component", "main").Logger()

	lg.Info().Msg("delayed-notifier application started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.New()
	if err != nil {
		lg.Error().Err(err).Msg("application initialization failed")
		wait()
		os.Exit(1)
	}

	if err := app.Run(cancel); err != nil {
		lg.Error().Err(err).Msg("application stopped with error")
		wait()
		os.Exit(1)
	}

	app.GracefullShutdown(ctx, cancel)
	lg.Info().Msg("application stopped gracefully")
	wait()
}

func wait() {
	fmt.Println("Press Enter to close…")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
