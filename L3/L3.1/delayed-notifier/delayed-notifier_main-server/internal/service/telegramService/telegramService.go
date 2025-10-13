package telegramService

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/telegram"
	"github.com/wb-go/wbf/zlog"
)

type BotUpdate tgbotapi.Update

type TelegramService struct {
	lg zlog.Zerolog
	tg *telegram.Client
}

func New(tg *telegram.Client) *TelegramService {
	lg := zlog.Logger.With().Str("component", "service-telegramService").Logger()
	return &TelegramService{
		lg: lg,
		tg: tg,
	}
}

func (sv *TelegramService) HandleStart(chatID int64, message string) error {
	return sv.tg.HandleStart(chatID, message)
}
