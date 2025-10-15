package service

import (
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	consumeNoticeServiceConfig *consumeNoticeService.Config
}

type retryStrategy struct {
	Attempts int
	Delay    time.Duration
	Backoff  float64
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		consumeNoticeServiceConfig: consumeNoticeService.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("service:\n%s", c.consumeNoticeServiceConfig.String())
}
