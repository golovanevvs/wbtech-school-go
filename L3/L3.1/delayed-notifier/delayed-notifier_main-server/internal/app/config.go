package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/logger"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	lg *logger.Config
	tr *transport.Config
	rp *repository.Config
}

func newConfig() (*appConfig, error) {
	appConfig := &appConfig{
		tr: transport.NewConfig(),
		rp: repository.NewConfig(),
		lg: logger.NewConfig(),
	}

	envFilePath := ".env"

	appConfigFilePath := "./infra/app/config.yaml"
	postgresConfigFilePath := "./infra/postgres/config.yaml"

	cfg := config.New()

	// -----------------------------------
	// Define flags
	// -----------------------------------

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

	if err := cfg.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// -----------------------------------
	// Load config for app
	// -----------------------------------

	// ------------- Logger -------------
	// appConfig.lg.Level = cfg.GetString("logger.level")

	// ------------- Transport -------------
	// appConfig.tr.TrHTTP.Port = cfg.GetInt("server.port")
	// appConfig.tr.TrHTTP.Handler.GinMode = cfg.GetString("gin.mode")

	// -----------------------------------
	// Load config for postgres
	// -----------------------------------

	// appConfig.rp.Postgres.Master.Host = cfg.GetString("postgres.host")
	// appConfig.rp.Postgres.Master.Port = cfg.GetInt("postgres.master.port")
	// appConfig.rp.Postgres.Master.User = cfg.GetString("postgres.user")
	// appConfig.rp.Postgres.Master.Password = cfg.GetString("postgres.password")
	// appConfig.rp.Postgres.Master.DBName = cfg.GetString("postgres.db")

	// appConfig.rp.Postgres.Slave1.Host = cfg.GetString("postgres.host")
	// appConfig.rp.Postgres.Slave1.Port = cfg.GetInt("postgres.slave1.port")
	// appConfig.rp.Postgres.Slave1.User = cfg.GetString("postgres.user")
	// appConfig.rp.Postgres.Slave1.Password = cfg.GetString("postgres.password")
	// appConfig.rp.Postgres.Slave1.DBName = cfg.GetString("postgres.db")

	// appConfig.rp.Postgres.Slave2.Host = cfg.GetString("postgres.host")
	// appConfig.rp.Postgres.Slave2.Port = cfg.GetInt("postgres.slave2.port")
	// appConfig.rp.Postgres.Slave2.User = cfg.GetString("postgres.user")
	// appConfig.rp.Postgres.Slave2.Password = cfg.GetString("postgres.password")
	// appConfig.rp.Postgres.Slave2.DBName = cfg.GetString("postgres.db")

	// appConfig.rp.Postgres.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	// appConfig.rp.Postgres.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	// appConfig.rp.Postgres.ConnMaxLifetime = cfg.GetDuration("postgres.conn_max_lifetime")

	return appConfig, nil
}
