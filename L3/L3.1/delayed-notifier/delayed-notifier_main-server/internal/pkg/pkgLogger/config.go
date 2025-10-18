package pkgLogger

import (
	"fmt"

	"github.com/fatih/color"
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
	return fmt.Sprintf(`logger:
  %s: %s`,
		color.YellowString("level"), color.GreenString(c.Level))
}
