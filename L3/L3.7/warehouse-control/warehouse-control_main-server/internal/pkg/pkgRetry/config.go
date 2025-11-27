package pkgRetry

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Attempts int
	Delay    time.Duration
	Backoff  float64
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Attempts: cfg.GetInt("app.retry.attempts"),
		Delay:    cfg.GetDuration("app.retry.delay"),
		Backoff:  cfg.GetFloat64("app.retry.backoff"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`retry strategy:
  %s: %d, %s: %v, %s: %v`,
		"attempts", c.Attempts,
		"delay", c.Delay,
		"backoff", c.Backoff)
}
