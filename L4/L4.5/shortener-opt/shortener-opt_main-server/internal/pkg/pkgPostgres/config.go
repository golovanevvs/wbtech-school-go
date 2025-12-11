package pkgPostgres

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Master          *commonConfig
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type commonConfig struct {
	Port int
	Host string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Master: &commonConfig{
			Port: cfg.GetInt("postgres.master.port"),
			Host: cfg.GetString("postgres.master.host"),
		},
		User:            cfg.GetString("postgres.user"),
		Password:        cfg.GetString("postgres.password"),
		DBName:          cfg.GetString("postgres.db"),
		SSLMode:         cfg.GetString("postgres.sslmode"),
		MaxOpenConns:    cfg.GetInt("postgres.max_open_conns"),
		MaxIdleConns:    cfg.GetInt("postgres.max_idle_conns"),
		ConnMaxLifetime: cfg.GetDuration("postgres.conn_max_lifetime"),
	}
}

func (c commonConfig) String() string {
	return fmt.Sprintf(`%s: %s, %s: %d`,
		"Host", c.Host, "Port", c.Port)
}

func (c Config) String() string {
	var password string
	if c.Password != "" {
		password = "***"
	}
	return fmt.Sprintf(`postgres:
  %s: %s
  %s: %s, %s: %s, %s: %s, %s:%s
  %s: %d, %s: %d, %s: %v`,
		"Master", c.Master.String(),
		"User", c.User, "Password", password, "DBName", c.DBName, "SSL mode", c.SSLMode,
		"MaxOpenConns", c.MaxOpenConns, "MaxIdleConns", c.MaxIdleConns, "ConnMaxLifetime", c.ConnMaxLifetime)
}
