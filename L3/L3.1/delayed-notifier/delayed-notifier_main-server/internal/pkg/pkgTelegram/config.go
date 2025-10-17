package pkgTelegram

import (
	"github.com/wb-go/wbf/config"
)

type Config struct {
	Token string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Token: cfg.GetString("telegram.token"),
	}
}

func (c Config) String() string {
	return "telegram:\n \033[33mtoken: \033[0m\033[32m***\033[0m"
}
