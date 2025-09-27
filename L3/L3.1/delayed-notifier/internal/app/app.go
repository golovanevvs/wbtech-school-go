package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wb-go/wbf/zlog"
)

func Run() {
	zlog.InitConsole()

	log := zlog.Logger.With().Str("component", "app").Logger()

	log.Info().Msg("delayed-notifier app started")

	appConfig, err := NewAppConfig("./config/config.yaml", "./.env", "")
	if err != nil {
		log.Error().Err(err).Msg("error creating configuration")
	}

	fmt.Printf("addr: %s\n", appConfig.serverConfig.addr)
	fmt.Printf("logLevel: %s\n", appConfig.loggerConfig.logLevel)

	fmt.Println("Press Enter to exit…")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
