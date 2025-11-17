package pkgKafka

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Brokers []string
	Topic   string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Brokers: cfg.GetStringSlice("kafka.brokers"),
		Topic:   cfg.GetString("kafka.topic"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`kafka:
  %s: %v,
  %s: %s`,
		"brokers", c.Brokers,
		"topic", c.Topic,
	)
}
