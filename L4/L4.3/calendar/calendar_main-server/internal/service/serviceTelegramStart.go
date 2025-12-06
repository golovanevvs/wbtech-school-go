package service

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
)

// ITelegramHandler interface for telegram bot handlers
type ITelegramHandler interface {
	HandleCommand(ctx context.Context, chatID int64, command, message string) error
}

// TelegramHandler handles basic Telegram bot commands for calendar
type TelegramHandler struct {
	tg *pkgTelegram.Client
}

// NewTelegramHandler creates a new TelegramHandler
func NewTelegramHandler(tg *pkgTelegram.Client) *TelegramHandler {
	return &TelegramHandler{
		tg: tg,
	}
}

// HandleCommand handles incoming Telegram commands
func (th *TelegramHandler) HandleCommand(ctx context.Context, chatID int64, command, message string) error {
	switch command {
	case "/start":
		return th.handleStart(ctx, chatID, message)
	case "/help":
		return th.handleHelp(ctx, chatID)
	case "/events":
		return th.handleEvents(ctx, chatID)
	default:
		return th.handleUnknown(ctx, chatID, command)
	}
}

// handleStart handles the /start command
func (th *TelegramHandler) handleStart(ctx context.Context, chatID int64, message string) error {
	welcomeMessage := `🤖 Добро пожаловать в Calendar Bot!

Я помогу вам управлять событиями в календаре.

Доступные команды:
/help - показать справку
/events - получить ближайшие события

Просто отправьте мне сообщение, и я помогу!`

	return th.tg.SendTo(chatID, welcomeMessage)
}

// handleHelp handles the /help command
func (th *TelegramHandler) handleHelp(ctx context.Context, chatID int64) error {
	helpMessage := `📚 Справка по командам:

/start - приветствие и начало работы
/help - показать эту справку
/events - показать ближайшие события

💡 События создаются через веб-интерфейс календаря, а я могу отправлять напоминания!`

	return th.tg.SendTo(chatID, helpMessage)
}

// handleEvents handles the /events command
func (th *TelegramHandler) handleEvents(ctx context.Context, chatID int64) error {
	eventsMessage := `📅 Ближайшие события:

В данный момент события отображаются в веб-интерфейсе календаря.
Для получения напоминаний убедитесь, что события созданы с включенным напоминанием.

Создавайте события на сайте: / (главная страница)`

	return th.tg.SendTo(chatID, eventsMessage)
}

// handleUnknown handles unknown commands
func (th *TelegramHandler) handleUnknown(ctx context.Context, chatID int64, command string) error {
	unknownMessage := fmt.Sprintf(`❓ Неизвестная команда: %s

Используйте /help для просмотра доступных команд.`, command)

	return th.tg.SendTo(chatID, unknownMessage)
}

// SendTestMessage sends a test message to verify bot functionality
func (th *TelegramHandler) SendTestMessage(ctx context.Context, chatID int64) error {
	testMessage := "🧪 Тестовое сообщение от Calendar Bot"

	return th.tg.SendTo(chatID, testMessage)
}
