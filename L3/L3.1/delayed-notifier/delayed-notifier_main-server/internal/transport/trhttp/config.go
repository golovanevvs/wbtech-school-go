package trhttp

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	Port                       int
	PublicHost                 string
	RetryStrategyForWaitServer RetryStrategy
	Handler                    *handler.Config
}

type RetryStrategy struct {
	Attempts int
	Delay    time.Duration
	Backoff  float64
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Port:       cfg.GetInt("app.transport.http.port"),
		PublicHost: cfg.GetString("app.transport.http.public_host"),
		RetryStrategyForWaitServer: RetryStrategy{
			Attempts: cfg.GetInt("app.transport.http.retry-strategy.attempts"),
			Delay:    cfg.GetDuration("app.transport.http.retry-strategy.delay"),
			Backoff:  cfg.GetFloat64("app.transport.http.retry-strategy.backoff"),
		},
		Handler: handler.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`http:
    %s: %s
    %s: %s
    retry strategy for wait server:
      %s: %s, %s: %s, %s: %s
    %s`,
		color.YellowString("public host"), color.GreenString(c.PublicHost),
		color.YellowString("port"), color.GreenString("%d", c.Port),
		color.YellowString("attempts"), color.GreenString("%d", c.RetryStrategyForWaitServer.Attempts),
		color.YellowString("delay"), color.GreenString("%v", c.RetryStrategyForWaitServer.Delay),
		color.YellowString("backoff"), color.GreenString("%v", c.RetryStrategyForWaitServer.Backoff),
		c.Handler.String(),
	)
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Port)
	}

	return nil
}
