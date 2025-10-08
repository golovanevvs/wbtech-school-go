package repository

import "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/sql/postgres"

type Config struct {
	Postgres *postgres.Config
}

func NewConfig() *Config {
	return &Config{
		Postgres: postgres.NewConfig(),
	}
}
