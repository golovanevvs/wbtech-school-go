package email

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	SMTPPort int
	SMTPHost string
	Username string
	Password string
	From     string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		SMTPPort: cfg.GetInt("email.smtp_port"),
		SMTPHost: cfg.GetString("email_smtp_host"),
		Username: cfg.GetString("email.username"),
		Password: cfg.GetString("email.password"),
		From:     cfg.GetString("email.from"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("email:\n \033[33mSMTPHost: \033[0m\033[32m%s\033[0m, \033[33mSMTPPort: \033[0m\033[32m%d\033[0m\n\033[33musername: \033[0m\033[32m%s\033[0m, \033[33mfrom: \033[0m\033[32m%s\033[0m", c.SMTPHost, c.SMTPPort, c.Username, c.From)
}
