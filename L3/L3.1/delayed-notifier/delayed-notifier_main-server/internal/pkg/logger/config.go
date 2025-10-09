package logger

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Level string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Level: cfg.GetString("logger.level"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("logger:\n \033[33mlevel: \033[0m\033[32m%s\033[0m", c.Level)
}
