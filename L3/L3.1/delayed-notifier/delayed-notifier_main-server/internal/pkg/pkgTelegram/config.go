package pkgTelegram

import (
	"fmt"

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
	return fmt.Sprintf("telegram:\n \033[33mtoken: \033[0m\033[32m%s\033[0m", c.Token)
}
