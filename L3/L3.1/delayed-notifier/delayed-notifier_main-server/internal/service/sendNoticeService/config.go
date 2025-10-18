package sendNoticeService

import (
	"fmt"
	"time"

	"github.com/fatih/color"
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
			Attempts: cfg.GetInt("app.service.sendnoticeservice.retry-strategy.attempts"),
			Delay:    cfg.GetDuration("app.service.sendnoticeservice.retry-strategy.delay"),
			Backoff:  cfg.GetFloat64("app.service.sendnoticeservice.retry-strategy.backoff"),
		},
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`SendNoticeService:
    retry strategy:
      %s: %s, %s: %s, %s: %s`,
		color.YellowString("attempts"), color.GreenString("%d", c.RetryStrategy.Attempts),
		color.YellowString("delay"), color.GreenString("%v", c.RetryStrategy.Delay),
		color.YellowString("backoff"), color.GreenString("%v", c.RetryStrategy.Backoff))
}
