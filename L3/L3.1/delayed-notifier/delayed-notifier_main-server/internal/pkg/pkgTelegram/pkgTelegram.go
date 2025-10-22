package pkgTelegram

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wb-go/wbf/zlog"
)

type Client struct {
	lg  zlog.Zerolog
	bot *tgbotapi.BotAPI
}

func New(cfg *Config) (*Client, error) {
	lg := zlog.Logger.With().Str("component", "telegram").Logger()

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	lg.Debug().Str("account", bot.Self.UserName).Msg("authorized on account")

	return &Client{
		lg:  lg,
		bot: bot}, nil
}

func (c *Client) SendTo(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	c.lg.Debug().Int64("chatID", chatID).Str("message", message).Msg("message sent")
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
	c.lg.Debug().Msg("checking telegram webhook...")

	info, err := c.bot.GetWebhookInfo()
	if err != nil {
		c.lg.Error().Err(err).Msg("cannot get telegram webhook info")
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	c.lg.Debug().Str("URL", info.URL).Bool("has custom certificate", info.HasCustomCertificate).Int("pending update count", info.PendingUpdateCount).Msg("current webhook info")

	if info.URL == url {
		c.lg.Info().Msg("the webhook has not changed")
		return nil
	}

	c.lg.Debug().Msg("setting telegram webhook...")

	webhookConfig, err := tgbotapi.NewWebhook(url)

	if err != nil {
		return fmt.Errorf("failed to create webhook config: %w", err)
	}

	_, err = c.bot.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	time.Sleep(1 * time.Second)

	info, err = c.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	c.lg.Debug().Str("webhook", info.URL).Msg("webhook set")

	return nil
}

func (c *Client) HandleStart(chatID int64, message string) error {
	if strings.HasPrefix(message, "/start") {
		if err := c.SendTo(chatID, "Telegram successfully linked!"); err != nil {
			return err
		}
		c.lg.Debug().Int64("chatID", chatID).Msg("Telegram successfully linked")
		return nil
	}

	if err := c.SendTo(chatID, "I only understand the command /start."); err != nil {
		return err
	}

	return fmt.Errorf("unknown command")
}
