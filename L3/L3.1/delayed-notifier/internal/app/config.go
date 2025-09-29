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
	port int
}

type loggerConfig struct {
	logLevel string
}

func newAppConfig(configFilePath, envFilePath, envPrefix string) (*appConfig, error) {
	appConfig := &appConfig{}

	cfg := config.New()

	cfg.DefineFlag("a", "addr", "server.addr", ":7777", "Server address")

	cfg.ParseFlags()

	err := cfg.Load(configFilePath, envFilePath, envPrefix)
	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig.serverConfig.port = cfg.GetInt("server.port")
	appConfig.loggerConfig.logLevel = cfg.GetString("logger.level")

	err = appConfig.validate()
	if err != nil {
		return appConfig, fmt.Errorf("failed to validation config: %w", err)
	}

	return appConfig, nil
}

func (ap *appConfig) validate() error {
	if ap.serverConfig.port <= 0 || ap.serverConfig.port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", ap.serverConfig.port)
	}

	return nil
}
