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
	return fmt.Sprintf("redis:\n \033[33mHost:\033[0m \033[32m%s\033[0m, \033[33mPort:\033[0m \033[32m%d\033[0m, \033[33mDB:\033[0m \033[32m%d\033[0m, \033[33mTTL:\033[0m \033[32m%v\033[0m", c.Host, c.Port, c.DB, c.TTL)
}
