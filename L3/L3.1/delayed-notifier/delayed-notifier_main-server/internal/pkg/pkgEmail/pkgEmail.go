package pkgEmail

import (
	"fmt"
	"net/smtp"
)

type Client struct {
	cfg Config
}

func New(cfg Config) (*Client, error) {
	if cfg.SMTPHost == "" || cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("invalid email config")
	}
	return &Client{cfg: cfg}, nil
}

func (c *Client) Send(to string, message string) error {
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", c.cfg.SMTPHost, c.cfg.SMTPPort)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Notification\r\n"+
			"\r\n"+
			"%s\r\n", to, message))

	if err := smtp.SendMail(addr, auth, c.cfg.From, []string{to}, msg); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
