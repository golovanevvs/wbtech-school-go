package service

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/sendNoticeService"
	"github.com/wb-go/wbf/config"
)

type Config struct {
	sendNoticeServiceConfig *sendNoticeService.Config
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		sendNoticeServiceConfig: sendNoticeService.NewConfig(cfg),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`service:
  %s`, c.sendNoticeServiceConfig.String())
}
