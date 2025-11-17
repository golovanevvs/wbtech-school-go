package repository

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/repository/rpFileStorage"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	RpFileStorage *rpFileStorage.Config
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		rpFileStorage.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`repository:
  %s`,
		c.RpFileStorage.String())
}
