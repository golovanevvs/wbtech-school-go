package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/email"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/logger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/telegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	lg *logger.Config
	tr *transport.Config
	rp *repository.Config
	rb *rabbitmq.Config
	tg *telegram.Config
	em *email.Config
}

func newConfig() (*Config, error) {

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

	cfg.DefineFlag("p", "srvport", "transport.http.port", 6000, "HTTP server port")
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	return &Config{
		lg: logger.NewConfig(cfg),
		tr: transport.NewConfig(cfg),
		rp: repository.NewConfig(cfg),
		rb: rabbitmq.NewConfig(cfg),
		tg: telegram.NewConfig(cfg),
		em: email.NewConfig(cfg),
	}, nil
}

func (a *Config) String() string {
	if a == nil {
		return "appConfig: <nil>"
	}
	return fmt.Sprintf("Configuration:\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n",
		a.lg.String(),
		a.tr.String(),
		a.rp.String(),
		a.rb.String(),
		a.tg.String(),
		a.em.String(),
	)
}
