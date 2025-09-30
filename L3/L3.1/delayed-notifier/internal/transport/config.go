package transport

import transporthttp "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/transport/http"

type Config struct {
	TrHTTP *transporthttp.Config
}

func NewConfig() *Config {
	return &Config{
		TrHTTP: transporthttp.NewConfig(),
	}
}
