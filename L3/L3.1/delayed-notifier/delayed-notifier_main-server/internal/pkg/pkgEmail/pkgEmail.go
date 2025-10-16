package pkgEmail

import (
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"
)

type Client struct {
	cfg  *Config
	lock sync.Mutex
}

func New(cfg *Config) (*Client, error) {
	if cfg.SMTPHost == "" || cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("invalid email config")
	}
	return &Client{cfg: cfg}, nil
}

func (c *Client) SendEmail(to []string, subject, body string, isHTML bool, attachments ...string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	msg := gomail.NewMessage()
	msg.SetHeader("From", c.cfg.From)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)

	if isHTML {
		msg.SetBody("text/html", body)
	} else {
		msg.SetBody("text/plain", body)
	}

	for _, a := range attachments {
		msg.Attach(a)
	}

	dialer := gomail.NewDialer(c.cfg.SMTPHost, c.cfg.SMTPPort, c.cfg.Username, c.cfg.Password)
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
