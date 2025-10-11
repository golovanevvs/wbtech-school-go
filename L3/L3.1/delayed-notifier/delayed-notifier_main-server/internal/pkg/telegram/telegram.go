package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

type Message struct {
	ChatID int64
	Text   string
}

func New(cfg *Config) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &Client{bot: bot}, nil
}

func (c *Client) SendTo(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Printf("Sent message to chat %d: %s", chatID, message)
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

func (c *Client) ListenUpdates(handler func(update tgbotapi.Update)) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)
	for update := range updates {
		handler(update)
	}
}

/*
tg, _ := telegram.New(&telegram.Config{Token: "YOUR_BOT_TOKEN"})

tg.ListenUpdates(func(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() && update.Message.Command() == "start" {
		args := update.Message.CommandArguments() // всё, что после /start
		log.Printf("User %s started bot with args: %s", update.Message.From.UserName, args)

		chatID := update.Message.Chat.ID
		// сохранить chatID в БД и привязать к args (например user_id)
		tg.SendTo(chatID, "✅ Telegram успешно привязан!")
	}
})
*/

func (c *Client) SetWebhook(url string) error {
	webhookConfig, err := tgbotapi.NewWebhook(url)
	if err != nil {
		return fmt.Errorf("failed to create webhook config: %w", err)
	}

	_, err = c.bot.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	info, err := c.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	log.Printf("Webhook set to: %s", info.URL)
	return nil
}

func (c *Client) HandleUpdate(update tgbotapi.Update) *Message {
	if update.Message == nil {
		return nil
	}

	return &Message{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	}
}

func (c *Client) HandleStart(msg *Message) bool {
	if msg == nil {
		return false
	}

	if msg.Text == "/start" || len(msg.Text) > 6 && msg.Text[:6] == "/start" {
		log.Printf("User %d started bot with args: %s", msg.ChatID, msg.Text[6:])
		c.SendTo(msg.ChatID, "✅ Telegram успешно привязан!")
		return true
	}

	return false
}

/*
r.POST("/telegram/webhook", func(c *gin.Context) {
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.Status(400)
		return
	}

	msg := tg.HandleUpdate(update)
	if msg != nil {
		if !tg.HandleStart(msg) {
			// если это не /start — можно обработать иначе
			tg.SendTo(msg.ChatID, "Я понимаю только команду /start 🙂")
		}
	}

	c.Status(200)
})

*/
