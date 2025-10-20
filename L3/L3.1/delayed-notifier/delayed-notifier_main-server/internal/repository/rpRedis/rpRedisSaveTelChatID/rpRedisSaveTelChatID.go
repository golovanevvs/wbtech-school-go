package rpRedisSaveTelChatID

import (
	"context"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisSaveChatID struct {
	lg *zlog.Zerolog
	rd *pkgRedis.Client
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisSaveChatID {
	lg := parentLg.With().Str("component-2", "RpRedisSaveTelChatID").Logger()
	return &RpRedisSaveChatID{
		lg: &lg,
		rd: rd,
	}
}

func (rp *RpRedisSaveChatID) SaveTelegramChatID(ctx context.Context, username string, chatID int64) (err error) {
	lg := rp.lg.With().Str("method", "SaveTelegramChatID").Logger()
	lg.Trace().Msg("⬇ method starting")
	defer lg.Trace().Msg("⬆ method stopped")

	lg.Trace().Str("username", username).Int64("chat ID", chatID).Msgf("%s saving name, chat ID to Redis...", color.YellowString("➤"))
	err = rp.rd.Set(ctx, username, chatID, 0)
	if err != nil {
		return pkgErrors.Wrapf(err, "save to Redis, name: %s, chat ID: %d", username, chatID)
	}
	lg.Trace().Str("username", username).Int64("chat ID", chatID).Msgf("%s name, chat ID saved to Redis successfully", color.GreenString("✔"))

	return nil
}
