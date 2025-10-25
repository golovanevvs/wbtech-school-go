package pkgEmail

import (
	"crypto/tls"
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"
)

type Dialer interface {
	DialAndSend(...*gomail.Message) error
}

var newDialer = func(host string, port int, user, pass string) Dialer {
	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

type Client struct {
	cfg  *Config
	lock sync.Mutex
}

func New(cfg *Config) (*Client, error) {
	if cfg.SMTPHost == "" || cfg.Username == "" || cfg.Password == "" || cfg.From == "" {
		return nil, fmt.Errorf("invalid email config: missing required fields")
	}
	return &Client{
		cfg: cfg,
	}, nil
}

func (c *Client) SendEmail(to []string, subject, body string, isHTML bool, attachments ...string) error {

	msg := gomail.NewMessage()
	c.lock.Lock()
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
	c.lock.Unlock()

	dialer := newDialer(c.cfg.SMTPHost, c.cfg.SMTPPort, c.cfg.Username, c.cfg.Password)
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
