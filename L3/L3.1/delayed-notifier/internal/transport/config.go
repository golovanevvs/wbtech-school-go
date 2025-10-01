package transport

import "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/transport/trhttp"

type Config struct {
	TrHTTP *trhttp.Config
}

func NewConfig() *Config {
	return &Config{
		TrHTTP: trhttp.NewConfig(),
	}
}
