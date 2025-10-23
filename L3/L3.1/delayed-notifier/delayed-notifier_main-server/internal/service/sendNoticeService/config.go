package sendNoticeService

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	RetryStrategy retryStrategy
}

type retryStrategy struct {
	Attempts int
	Delay    time.Duration
	Backoff  float64
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		RetryStrategy: retryStrategy{
			Attempts: cfg.GetInt("app.service.sendnoticeservice.retry.attempts"),
			Delay:    cfg.GetDuration("app.service.sendnoticeservice.retry.delay"),
			Backoff:  cfg.GetFloat64("app.service.sendnoticeservice.retry.backoff"),
		},
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`SendNoticeService:
    %s:
      %s: %d, %s: %v, %s: %v`,
		"retry strategy",
		"attempts", c.RetryStrategy.Attempts,
		"delay", c.RetryStrategy.Delay,
		"backoff", c.RetryStrategy.Backoff)
}
