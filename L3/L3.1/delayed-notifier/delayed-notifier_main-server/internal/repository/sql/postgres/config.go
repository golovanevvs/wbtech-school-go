package postgres

import "time"

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

func NewConfig() *Config {
	return &Config{
		Master: &dsnConfig{},
		Slave1: &dsnConfig{},
		Slave2: &dsnConfig{},
	}
}
