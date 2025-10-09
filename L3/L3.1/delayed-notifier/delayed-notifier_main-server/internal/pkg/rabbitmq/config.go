package rabbitmq

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
	return fmt.Sprintf("rabbitmq:\n \033[33mHost:\033[0m \033[32m%s\033[0m, \033[33mPort:\033[0m \033[32m%d\033[0m, \033[33mVHost:\033[0m \033[32m%s\033[0m, \033[33mUsername:\033[0m \033[32m%s\033[0m\n \033[33mExchange:\033[0m \033[32m%s\033[0m, \033[33mExchangeType:\033[0m \033[32m%s\033[0m\n \033[33mQueue:\033[0m \033[32m%s\033[0m, \033[33mRoutingKey:\033[0m \033[32m%s\033[0m\n \033[33mDLX:\033[0m \033[32m%s\033[0m, \033[33mDLQ:\033[0m \033[32m%s\033[0m", c.Host, c.Port, c.VHost, c.Username, c.Exchange, c.ExchangeType, c.Queue, c.RoutingKey, c.DLX, c.DLQ)
}
