package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/app"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()

	zlog.Logger.Info().Str("component", "main").Msg("delayed-notifier application started")

	app, err := app.New()
	if err != nil {
		wait()
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		wait()
		os.Exit(1)
	}

	wait()
	zlog.Logger.Info().Str("component", "main").Msg("delayed-notifier application stopped")
}

func wait() {
	fmt.Println("Press Enter to exit…")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
