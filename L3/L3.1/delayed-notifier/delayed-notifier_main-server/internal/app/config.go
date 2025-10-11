package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/logger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	lg *logger.Config
	tr *transport.Config
	rp *repository.Config
	rb *rabbitmq.Config
}

func newConfig() (*appConfig, error) {

	envFilePath := ".env"
	appConfigFilePath := "./providers/app/config-example.yaml"
	postgresConfigFilePath := "./providers/postgres/config-example.yaml"

	cfg := config.New()

	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	cfg.EnableEnv("")

	if err := cfg.LoadConfigFiles(appConfigFilePath, postgresConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	cfg.DefineFlag("p", "srvport", "transport.http.port", 7777, "HTTP server port")
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	appConfig := &appConfig{
		lg: logger.NewConfig(cfg),
		tr: transport.NewConfig(cfg),
		rp: repository.NewConfig(cfg),
		rb: rabbitmq.NewConfig(cfg),
	}

	return appConfig, nil
}

func (a *appConfig) String() string {
	if a == nil {
		return "appConfig: <nil>"
	}
	return fmt.Sprintf("Configuration:\n%s\n%s\n%s\n%s\n",
		a.lg.String(),
		a.tr.String(),
		a.rp.String(),
		a.rb.String())
}
