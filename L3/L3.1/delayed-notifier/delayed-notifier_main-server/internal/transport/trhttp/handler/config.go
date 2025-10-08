package handler

type Config struct {
	GinMode string `mapstructure:"gin.mode"`
}

func NewConfig() *Config {
	return &Config{}
}
