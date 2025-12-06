package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgLogger"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	lg *pkgLogger.Config
	rs *pkgRetry.Config
	tg *pkgTelegram.Config
	em *pkgEmail.Config
	pg *pkgPostgres.Config
	sv *service.Config
	tr *transport.Config
}

func newConfig(env string) (*Config, error) {
	cfg := config.New()

	if env == "local" {
		if err := cfg.LoadEnvFiles(
			".env",
			"providers/app/.env",
			"providers/email/.env",
			"providers/telegram/.env",
			"providers/postgres/.env",
		); err != nil {
			return nil, fmt.Errorf("failed to load env files: %w", err)
		}
	}

	cfg.EnableEnv("")

	if err := cfg.LoadConfigFiles(
		"providers/app/config.yaml",
		"providers/logger/config.yaml",
		"providers/postgres/config.yaml",
	); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	cfg.DefineFlag("p", "srvport", "app.transport.http.port", 6000, "HTTP server port")
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	return &Config{
		lg: pkgLogger.NewConfig(cfg),
		rs: pkgRetry.NewConfig(cfg),
		tg: pkgTelegram.NewConfig(cfg),
		em: pkgEmail.NewConfig(cfg),
		pg: pkgPostgres.NewConfig(cfg),
		sv: service.NewConfig(cfg),
		tr: transport.NewConfig(cfg, env),
	}, nil
}

func (a *Config) String() string {
	if a == nil {
		return "appConfig: <nil>"
	}
	return fmt.Sprintf(`%s Configuration:
%s %s

%s %s

%s %s

%s %s

%s %s

%s %s

%s %s
`,
		pkgConst.Config,
		pkgConst.Logger, a.lg.String(),
		pkgConst.Retry, a.rs.String(),
		pkgConst.Telegram, a.tg.String(),
		pkgConst.EMail, a.em.String(),
		pkgConst.Postgres, a.pg.String(),
		pkgConst.Service, a.sv.String(),
		pkgConst.Transport, a.tr.String(),
	)
}
