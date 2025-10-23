package pkgEmail

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
		SMTPHost: cfg.GetString("email.smtp_host"),
		Username: cfg.GetString("email.username"),
		Password: cfg.GetString("email.password"),
		From:     cfg.GetString("email.from"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`email:
  %s: %s, %s: %d
  %s: %s,
  %s: %s,`,
		"SMTPHost", c.SMTPHost, "SMTPPort", c.SMTPPort,
		"username", c.Username,
		"from", c.From)
}
