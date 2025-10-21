package pkgRabbitmq

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

// Config contains RabbitMQ connection and queue configuration.
type Config struct {
	Port         int
	Host         string
	VHost        string
	Username     string
	Password     string
	Exchange     string
	ExchangeType string
	Queue        string
	RoutingKey   string
	DLX          string
	DLQ          string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Port:         cfg.GetInt("rabbitmq.port"),
		Host:         cfg.GetString("rabbitmq.host"),
		VHost:        cfg.GetString("rabbitmq.vhost"),
		Username:     cfg.GetString("rabbitmq.username"),
		Password:     cfg.GetString("rabbitmq.password"),
		Exchange:     cfg.GetString("rabbitmq.exchange"),
		ExchangeType: cfg.GetString("rabbitmq.exchange_type"),
		Queue:        cfg.GetString("rabbitmq.queue"),
		RoutingKey:   cfg.GetString("rabbitmq.routing_key"),
		DLX:          cfg.GetString("rabbitmq.dlx"),
		DLQ:          cfg.GetString("rabbitmq.dlq"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("rabbitmq:\n Host: %s, Port: %d, VHost: %s,Username: %s\n Exchange: %s, ExchangeType: %s\n Queue:%s, RoutingKey: %s\n DLX: %s, DLQ: %s", c.Host, c.Port, c.VHost, c.Username, c.Exchange, c.ExchangeType, c.Queue, c.RoutingKey, c.DLX, c.DLQ)
}
