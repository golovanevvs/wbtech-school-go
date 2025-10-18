package transport

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	TrHTTP *trhttp.Config
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		TrHTTP: trhttp.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`transport:
  %s`,
		c.TrHTTP.String())
}
