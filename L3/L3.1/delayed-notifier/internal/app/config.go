package app

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	serverConfig serverConfig
	loggerConfig loggerConfig
}

type serverConfig struct {
	addr string
}

type loggerConfig struct {
	logLevel string
}

func NewAppConfig(configFilePath, envFilePath, envPrefix string) (*appConfig, error) {
	appConfig := &appConfig{}

	cfg := config.New()

	cfg.DefineFlag("a", "addr", "server.addr", ":7777", "Server address")

	cfg.ParseFlags()

	err := cfg.Load(configFilePath, envFilePath, envPrefix)
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig.serverConfig.addr = cfg.GetString("server.addr")
	appConfig.loggerConfig.logLevel = cfg.GetString("logger.level")

	return appConfig, nil
}
