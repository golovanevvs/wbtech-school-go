package telegramStartService

import (
	"context"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	SaveTelegramChatID(ctx context.Context, username string, chatID int64) (err error)
}

type TelegramStartService struct {
	lg *zlog.Zerolog
	tg *pkgTelegram.Client
	rp IRepository
}

func New(parentLg *zlog.Zerolog, tg *pkgTelegram.Client, rp IRepository) *TelegramStartService {
	lg := parentLg.With().Str("component", "TelegramStartService").Logger()
	return &TelegramStartService{
		lg: &lg,
		tg: tg,
		rp: rp,
	}
}

func (sv *TelegramStartService) Start(ctx context.Context, username string, chatID int64, message string) (err error) {
	lg := sv.lg.With().Str("method", "Start").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	lg.Trace().Str("username", username).Int64("chat ID", chatID).Str("message", message).Msgf("%s starting handle /start...", color.YellowString("➤"))
	err = sv.tg.HandleStart(chatID, message)
	if err != nil {
		return err
	}
	lg.Trace().Str("username", username).Int64("chat ID", chatID).Str("message", message).Msgf("%s Command /start executed successfully", color.GreenString("✔"))

	lg.Trace().Str("username", username).Int64("chat ID", chatID).Msgf("%s saving name, chat ID to repository...", color.YellowString("➤"))
	err = sv.rp.SaveTelegramChatID(ctx, username, chatID)
	if err != nil {
		return err
	}
	lg.Trace().Str("username", username).Int64("chat ID", chatID).Msgf("%s name, chat ID saved to repository successfully", color.GreenString("✔"))

	return nil
}
