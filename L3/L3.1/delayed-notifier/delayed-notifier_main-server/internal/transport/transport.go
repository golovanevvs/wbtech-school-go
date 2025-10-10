package transport

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp"
)

type IService interface {
	trhttp.IService
}

type Transport struct {
	HTTP *trhttp.HTTP
}

func New(cfg *Config, sv IService) *Transport {
	return &Transport{
		HTTP: trhttp.New(cfg.TrHTTP, sv),
	}
}
