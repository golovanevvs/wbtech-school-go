package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/logger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	Lg *logger.Config     `mapstructure:"logger"`
	Tr *transport.Config  `mapstructure:"transport"`
	Rp *repository.Config `mapstructure:"repository"`
}

func newConfig() (*appConfig, error) {

	envFilePath := ".env"

	appConfigFilePath := "./infra/app/config.yaml"
	postgresConfigFilePath := "./infra/postgres/config.yaml"

	cfg := config.New()

	if err := cfg.DefineFlag("p", "srvport", "server.port", 7777, "HTTP server port"); err != nil {
		return nil, fmt.Errorf("failed to define flags: %w", err)
	}

	cfg.ParseFlags()

	if err := cfg.LoadEnv(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env file: %w", err)
	}

	if err := cfg.Load(appConfigFilePath, ""); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if err := cfg.Load(postgresConfigFilePath, ""); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig := &appConfig{
		Lg: logger.NewConfig(cfg),
		Tr: transport.NewConfig(cfg),
		Rp: repository.NewConfig(cfg),
	}

	return appConfig, nil
}

func (a *appConfig) String() string {
	if a == nil {
		return "appConfig: <nil>"
	}
	return fmt.Sprintf("Configuration:\n%s\n%s\n%s", a.Lg.String(), a.Tr.String(), a.Rp.String())
}
