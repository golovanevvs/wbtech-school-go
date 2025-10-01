package handler

type Config struct {
	GinMode string
}

func NewConfig() *Config {
	return &Config{}
}
