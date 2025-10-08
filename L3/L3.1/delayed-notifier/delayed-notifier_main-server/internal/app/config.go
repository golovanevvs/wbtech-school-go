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

	err := cfg.loadEnv(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	err = cfg.Load(configFilePath, envPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// ------------- Logger -------------
	appConfig.lg.Level = cfg.GetString("logger.level")

	// ------------- Repository -------------
	appConfig.rp.Postgres.Master.Host = cfg.GetString("postgres.host")
	appConfig.rp.Postgres.Master.Port = cfg.GetInt("postgres.master.port")
	appConfig.rp.Postgres.Master.User = cfg.GetString("postgres.user")
	appConfig.rp.Postgres.Master.Password = cfg.GetString("postgres.password")
	appConfig.rp.Postgres.Master.DBName = cfg.GetString("postgres.db")

	appConfig.rp.Postgres.Slave1.Host = cfg.GetString("postgres.host")
	appConfig.rp.Postgres.Slave1.Port = cfg.GetInt("postgres.slave1.port")
	appConfig.rp.Postgres.Slave1.User = cfg.GetString("postgres.user")
	appConfig.rp.Postgres.Slave1.Password = cfg.GetString("postgres.password")
	appConfig.rp.Postgres.Slave1.DBName = cfg.GetString("postgres.db")

	appConfig.rp.Postgres.Slave2.Host = cfg.GetString("postgres.host")
	appConfig.rp.Postgres.Slave2.Port = cfg.GetInt("postgres.slave2.port")
	appConfig.rp.Postgres.Slave2.User = cfg.GetString("postgres.user")
	appConfig.rp.Postgres.Slave2.Password = cfg.GetString("postgres.password")
	appConfig.rp.Postgres.Slave2.DBName = cfg.GetString("postgres.db")

	appConfig.rp.Postgres.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.rp.Postgres.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.rp.Postgres.ConnMaxLifetime = cfg.GetDuration("postgres.conn_max_lifetime")

	// ------------- Transport -------------
	appConfig.tr.TrHTTP.Port = cfg.GetInt("server.port")
	appConfig.tr.TrHTTP.Handler.GinMode = cfg.GetString("gin.mode")

	return appConfig, nil
}
