package trhttp

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
)

type Config struct {
	Port    int `mapstructure:"port"`
	Handler *handler.Config
}

func NewConfig() *Config {
	return &Config{
		Handler: handler.NewConfig(),
	}
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Port)
	}

	return nil
}
