package logger

type Config struct {
	Level string
}

func NewConfig() *Config {
	return &Config{}
}
