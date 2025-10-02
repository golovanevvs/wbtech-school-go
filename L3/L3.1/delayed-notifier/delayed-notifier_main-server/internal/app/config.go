package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/pkg/logger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/transport"
	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	lg *logger.Config
	tr *transport.Config
	rp *repository.Config
}

func newConfig(configFilePath, envFilePath, envPrefix string) (*appConfig, error) {
	appConfig := &appConfig{
		tr: transport.NewConfig(),
		rp: repository.NewConfig(),
		lg: logger.NewConfig(),
	}

	cfg := config.New()

	if err := cfg.DefineFlag("p", "srvport", "server.port", "7777", "HTTP server port"); err != nil {
		return nil, fmt.Errorf("failed to define flags: %w", err)
	}

	cfg.ParseFlags()

	err := cfg.Load(configFilePath, envFilePath, envPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig.rp.Postgres.Master.Host = cfg.GetString("postgres.host")
	appConfig.tr.TrHTTP.Port = cfg.GetInt("server.port")
	appConfig.tr.TrHTTP.Handler.GinMode = cfg.GetString("gin.mode")
	appConfig.lg.Level = cfg.GetString("logger.level")

	return appConfig, nil
}
