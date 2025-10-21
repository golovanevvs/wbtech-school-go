package pkgRedis

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Port     int
	DB       int
	Host     string
	Password string
	TTL      time.Duration
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Port:     cfg.GetInt("redis.port"),
		DB:       cfg.GetInt("redis.db"),
		Host:     cfg.GetString("redis.host"),
		Password: cfg.GetString("redis.password"),
		TTL:      cfg.GetDuration("redis.ttl"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("redis:\n Host: %s, Port: %d, DB: %d, TTL: %v", c.Host, c.Port, c.DB, c.TTL)
}
