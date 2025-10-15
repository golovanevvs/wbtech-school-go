package consumeNoticeService

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
			Attempts: cfg.GetInt("app.service.consumenoticeservice.retry-strategy.attempts"),
			Delay:    cfg.GetDuration("app.service.consumenoticeservice.retry-strategy.delay"),
			Backoff:  cfg.GetFloat64("app.service.consumenoticeservice.retry-strategy.backoff"),
		},
	}
}

func (c Config) String() string {
	return fmt.Sprintf(" consumeNoticeService:\n  retry strategy:\n   attempts: \033[32m%d\033[0m delay: \033[32m%v\033[0m backoff: \033[32m%v\033[0m", c.RetryStrategy.Attempts, c.RetryStrategy.Delay, c.RetryStrategy.Backoff)
}
