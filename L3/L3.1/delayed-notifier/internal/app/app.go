package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wb-go/wbf/zlog"
)

func Run() {
	zlog.InitConsole()

	zlog.Logger.Info().Str("component", "app").Msg("delayed-notifier app started")

	appConfig, err := newAppConfig("./config/config.yaml", "./.env", "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "app").Msg("error creating configuration")
		os.Exit(1)
	}

	err = zlog.SetLevel(appConfig.loggerConfig.logLevel)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "logger").Msg("error set level")
	}

	// ------------------------- TEMP -------------------------
	fmt.Printf("addr: %s\n", appConfig.serverConfig.port)
	fmt.Printf("logLevel: %s\n", appConfig.loggerConfig.logLevel)

	fmt.Println("Press Enter to exit…")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
