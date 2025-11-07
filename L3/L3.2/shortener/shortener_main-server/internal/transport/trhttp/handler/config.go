package handler

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	GinMode       string
	WebClientPort int
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		GinMode:       cfg.GetString("app.transport.http.handler.gin_mode"),
		WebClientPort: cfg.GetInt("app.transport.http.handler.web_port"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`handler:
      %s: %s, %s: %d`,
		"Gin mode", c.GinMode, "Web client port", c.WebClientPort,
	)
}
