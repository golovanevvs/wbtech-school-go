package trhttp

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	Port       int
	PublicHost string
	Handler    *handler.Config
}

func NewConfig(cfg *config.Config, env string) *Config {
	var publicHost string
	if env == "local" {
		publicHost = cfg.GetString("app.transport.http.public_host")
	} else {
		publicHost = cfg.GetString("app.transport.http.public_host_docker")
	}
	return &Config{
		Port:       cfg.GetInt("app.transport.http.port"),
		PublicHost: publicHost,
		Handler:    handler.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`http:
    %s: %s
    %s: %d
    %s`,
		"public host", c.PublicHost,
		"port", c.Port,
		c.Handler.String(),
	)
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Port)
	}

	return nil
}
