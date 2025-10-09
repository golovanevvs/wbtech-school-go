package trhttp

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	Port    int
	Handler *handler.Config
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Port:    cfg.GetInt("transport.http.port"),
		Handler: handler.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(" http:\n\033[33m  port: \033[0m\033[32m%d\033[0m\n  handler:\n%s", c.Port, c.Handler.String())
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Port)
	}

	return nil
}
