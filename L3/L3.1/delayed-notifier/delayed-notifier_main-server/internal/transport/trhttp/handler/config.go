package handler

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	GinMode string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		GinMode: cfg.GetString("transport.http.handler.gin_mode"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("handler:\nGinMode: %s", c.GinMode)
}
