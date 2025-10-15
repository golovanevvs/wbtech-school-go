package telegramService

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/zlog"
)

type iRepository interface {
	SaveTelName(ctx context.Context, name string, chatID int64) (err error)
}

type TelegramService struct {
	lg zlog.Zerolog
	tg *pkgTelegram.Client
	rd *pkgRedis.Client
}

func New(tg *pkgTelegram.Client, rd *pkgRedis.Client) *TelegramService {
	lg := zlog.Logger.With().Str("component", "service-telegramService").Logger()
	return &TelegramService{
		lg: lg,
		tg: tg,
		rp: rp,
	}
}

func (sv *TelegramService) HandleStart(ctx context.Context, username string, chatID int64, message string) (err error) {
	err = sv.tg.HandleStart(chatID, message)
	if err != nil {
		return err
	}

	err = sv.rp.SaveTelName(ctx, username, chatID)
	if err != nil {
		return err
	}

	return nil
}
