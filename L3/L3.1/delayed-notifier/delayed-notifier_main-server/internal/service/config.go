package service

import noticeservice "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/noticeService"

type Config struct {
	NoticeService *noticeservice.Config
}

func NewConfig() *Config {
	return &Config{
		NoticeService: noticeservice.NewConfig(),
	}
}
