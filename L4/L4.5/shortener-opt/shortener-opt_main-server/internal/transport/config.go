package transport

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport/trhttp"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	TrHTTP *trhttp.Config
}

func NewConfig(cfg *config.Config, env string) *Config {
	return &Config{
		TrHTTP: trhttp.NewConfig(cfg, env),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`transport:
  %s`,
		c.TrHTTP.String())
}
