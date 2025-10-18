package handler

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	GinMode string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		GinMode: cfg.GetString("app.transport.http.handler.gin_mode"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`handler:
      %s: %s`,
		color.YellowString("Gin mode"), color.GreenString(c.GinMode),
	)
}
