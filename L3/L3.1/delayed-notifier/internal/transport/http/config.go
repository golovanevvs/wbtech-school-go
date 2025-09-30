package transporthttp

import "fmt"

type Config struct {
	Port int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Port)
	}

	return nil
}
