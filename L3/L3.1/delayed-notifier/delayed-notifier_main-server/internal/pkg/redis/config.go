package redis

import "time"

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}
