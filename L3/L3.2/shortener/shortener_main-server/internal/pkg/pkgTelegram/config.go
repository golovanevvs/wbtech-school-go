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
	var token string
	if len(c.Token) > 0 {
		token = "***hidden***"
	}
	return fmt.Sprintf(`telegram:
  %s: %s`,
		"token", token)
}
