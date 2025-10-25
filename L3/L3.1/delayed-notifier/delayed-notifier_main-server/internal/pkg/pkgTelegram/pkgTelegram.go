package pkgTelegram

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

func New(cfg *Config) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Client{bot: bot}, nil
}

func (c *Client) SendTo(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendToMany(chatIDs []int64, message string) error {
	var failed []int64
	for _, id := range chatIDs {
		if err := c.SendTo(id, message); err != nil {
			failed = append(failed, id)
		}
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to send to chatIDs: %v", failed)
	}
	return nil
}

func (c *Client) SetWebhook(url string) error {
	info, err := c.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("get webhook info: %w", err)
	}

	if info.URL == url {
		return nil
	}

	webhookConfig, err := tgbotapi.NewWebhook(url)
	if err != nil {
		return fmt.Errorf("create webhook config: %w", err)
	}

	_, err = c.bot.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	time.Sleep(1 * time.Second)

	info, err = c.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("get webhook info: %w", err)
	}

	return nil
}

func (c *Client) HandleStart(chatID int64, message string) error {
	if strings.HasPrefix(message, "/start") {
		if err := c.SendTo(chatID, "Telegram successfully linked!"); err != nil {
			return err
		}
		return nil
	}

	if err := c.SendTo(chatID, "I only understand the command /start."); err != nil {
		return err
	}

	return fmt.Errorf("unknown command")
}
