package logger

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Level string `mapstructure:"level"`
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Level: cfg.GetString("logger.level"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("\nlogger:\nlevel: %s", c.Level)
}
