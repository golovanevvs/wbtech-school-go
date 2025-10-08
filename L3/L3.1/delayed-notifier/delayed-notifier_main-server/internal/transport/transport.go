package transport

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp"
)

type Transport struct {
	HTTP *trhttp.HTTP
}

func New(cfg *Config) *Transport {
	return &Transport{
		HTTP: trhttp.New(cfg.TrHTTP),
	}
}
