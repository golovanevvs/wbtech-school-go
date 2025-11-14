package pkgPostgres

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Port     int
	Host     string
	User     string
	Password string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Port:     cfg.GetInt("POSTGRES_PORT"),
		Host:     cfg.GetString("POSTGRES_HOST"),
		User:     cfg.GetString("POSTGRES_USER"),
		Password: cfg.GetString("POSTGRES_PASSWORD"),
	}
}

func (c Config) String() string {
	var password string
	if c.Password != "" {
		password = "***"
	}
	return fmt.Sprintf(`postgres:
  %s: %s
  %s: %d
  %s: %s
  %s: %s`,
		"host", c.Host,
		"port", c.Port,
		"user", c.User,
		"password", password,
	)
}
