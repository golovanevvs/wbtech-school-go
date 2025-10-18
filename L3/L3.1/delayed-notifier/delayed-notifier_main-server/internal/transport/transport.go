package transport

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp"
	"github.com/wb-go/wbf/zlog"
)

type iService interface {
	trhttp.IService
}

type Transport struct {
	HTTP *trhttp.HTTP
}

func New(cfg *Config, sv iService) *Transport {
	lg := zlog.Logger.With().Str("layer", "transport").Logger()

	return &Transport{
		HTTP: trhttp.New(cfg.TrHTTP, &lg, sv),
	}
}
