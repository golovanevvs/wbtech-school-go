package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgLogger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	lg *pkgLogger.Config
	rd *pkgRedis.Config
	rb *pkgRabbitmq.Config
	tg *pkgTelegram.Config
	em *pkgEmail.Config
	sv *service.Config
	tr *transport.Config
}

func newConfig() (*Config, error) {

	envFilePath := ".env"
	appConfigFilePath := "./providers/app/config.yaml"
	redisConfigFilePath := "./providers/redis/config.yaml"
	rabbitmqConfigFilePath := "./providers/rabbitmq/config.yaml"
	postgresConfigFilePath := "./providers/postgres/config.yaml"

	cfg := config.New()

	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	cfg.EnableEnv("")

	if err := cfg.LoadConfigFiles(appConfigFilePath, redisConfigFilePath, rabbitmqConfigFilePath, postgresConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	cfg.DefineFlag("p", "srvport", "transport.http.port", 6000, "HTTP server port")
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	return &Config{
		lg: pkgLogger.NewConfig(cfg),
		rd: pkgRedis.NewConfig(cfg),
		rb: pkgRabbitmq.NewConfig(cfg),
		tg: pkgTelegram.NewConfig(cfg),
		em: pkgEmail.NewConfig(cfg),
		sv: service.NewConfig(cfg),
		tr: transport.NewConfig(cfg),
	}, nil
}

func (a *Config) String() string {
	if a == nil {
		return "appConfig: <nil>"
	}
	return fmt.Sprintf("Configuration:\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n",
		a.lg.String(),
		a.rd.String(),
		a.rb.String(),
		a.tg.String(),
		a.em.String(),
		a.sv.String(),
		a.tr.String(),
	)
}
