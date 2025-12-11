package main

import (
	"context"
	"os"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/app"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgConst"
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

	lg.Info().Msgf("%s shortener-opt_main-server started", pkgConst.AppStart)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.New(env)
	if err != nil {
		lg.Error().Err(err).Msgf("%s application initialization failed", pkgConst.Error)
		// wait()
		os.Exit(1)
	}

	if err := app.Run(cancel); err != nil {
		lg.Error().Err(err).Msgf("%s application stopped with error", pkgConst.Error)
		// wait()
		os.Exit(1)
	}

	app.GracefullShutdown(ctx, cancel)
	lg.Info().Msgf("%s application stopped gracefully", pkgConst.AppStop)
	// wait()
}

// func wait() {
// 	fmt.Println("Press Enter to close…")
// 	reader := bufio.NewReader(os.Stdin)
// 	_, _ = reader.ReadString('\n')
// }
