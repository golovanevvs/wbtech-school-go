package service

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

// Config service configuration for calendar
type Config struct {
	CleanupInterval       time.Duration
	ReminderCheckInterval time.Duration

	EnableTelegramNotifications bool
	EnableEmailNotifications    bool

	EventRetentionDays int
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		CleanupInterval:             cfg.GetDuration("app.service.cleanup_interval"),
		ReminderCheckInterval:       cfg.GetDuration("app.service.reminder_check_interval"),
		EnableTelegramNotifications: cfg.GetBool("app.service.enable_telegram_notifications"),
		EnableEmailNotifications:    cfg.GetBool("app.service.enable_email_notifications"),
		EventRetentionDays:          cfg.GetInt("app.service.event_retention_days"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`service:
  %s: %s
  %s: %s
  %s: %v
  %s: %v
  %s: %d days
  `,
		"cleanup_interval", c.CleanupInterval,
		"reminder_check_interval", c.ReminderCheckInterval,
		"enable_telegram_notifications", c.EnableTelegramNotifications,
		"enable_email_notifications", c.EnableEmailNotifications,
		"event_retention_days", c.EventRetentionDays,
	)
}
