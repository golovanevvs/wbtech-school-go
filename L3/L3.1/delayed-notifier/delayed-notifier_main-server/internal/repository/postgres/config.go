package postgres

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Master          *dsnConfig
	Slave1          *dsnConfig
	Slave2          *dsnConfig
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type dsnConfig struct {
	Port     int
	Host     string
	User     string
	Password string
	DBName   string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Master: &dsnConfig{
			Port:     cfg.GetInt("postgres.master.port"),
			Host:     cfg.GetString("postgres.host"),
			User:     cfg.GetString("postgres.user"),
			Password: cfg.GetString("postgres.password"),
			DBName:   cfg.GetString("postgres.db"),
		},
		Slave1: &dsnConfig{
			Port:     cfg.GetInt("postgres.slave1.port"),
			Host:     cfg.GetString("postgres.host"),
			User:     cfg.GetString("postgres.user"),
			Password: cfg.GetString("postgres.password"),
			DBName:   cfg.GetString("postgres.db"),
		},
		Slave2: &dsnConfig{
			Port:     cfg.GetInt("postgres.slave2.port"),
			Host:     cfg.GetString("postgres.host"),
			User:     cfg.GetString("postgres.user"),
			Password: cfg.GetString("postgres.password"),
			DBName:   cfg.GetString("postgres.db"),
		},
		MaxOpenConns:    cfg.GetInt("postgres.max_open_conns"),
		MaxIdleConns:    cfg.GetInt("postgres.max_idle_conns"),
		ConnMaxLifetime: cfg.GetDuration("postgres.conn_max_lifetime"),
	}
}

func (c dsnConfig) String() string {
	return fmt.Sprintf("\033[33mHost: \033[0m\033[32m%s\033[0m, \033[33mPort: \033[0m\033[32m%d\033[0m, \033[33mUser: \033[0m\033[32m%s\033[0m, \033[33mDBName: \033[0m\033[32m%s\033[0m",
		c.Host, c.Port, c.User, c.DBName)
}

func (c Config) String() string {
	return fmt.Sprintf("postgres:\n  Master: \033[32m%s\033[0m\n  Slave1: \033[32m%s\033[0m\n  Slave2: \033[32m%s\033[0m\n\033[33m  MaxOpenConns:\033[0m \033[32m%d\033[0m, \033[33mMaxIdleConns:\033[0m \033[32m%d\033[0m, \033[33mConnMaxLifetime:\033[0m \033[32m%s\033[0m\n", c.Master.String(), c.Slave1.String(), c.Slave2.String(), c.MaxOpenConns, c.MaxIdleConns, c.ConnMaxLifetime)
}
