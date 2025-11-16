package transport

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/transport/trhttp"
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
