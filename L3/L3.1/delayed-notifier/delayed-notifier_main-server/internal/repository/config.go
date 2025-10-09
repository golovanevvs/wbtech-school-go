package repository

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/postgres"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	Postgres *postgres.Config `mapstructure:"postgres"`
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Postgres: postgres.NewConfig(cfg),
	}
}

func (c *Config) String() string {
	if c == nil {
		return "Config: <nil>"
	}
	return fmt.Sprintf("\nrepository:\n %s", c.Postgres.String())
}
