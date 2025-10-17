package rpRedisLoadTelChatID

import (
	"context"
	"strconv"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisLoadTelChatID struct {
	lg *zlog.Zerolog
	rd *pkgRedis.Client
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisLoadTelChatID {
	lg := parentLg.With().Str("component1", "RpRedisLoadTelChatID").Logger()
	return &RpRedisLoadTelChatID{
		lg: &lg,
		rd: rd,
	}
}

func (rp *RpRedisLoadTelChatID) LoadTelegramChatID(ctx context.Context, username string) (chatID int64, err error) {
	lg := rp.lg.With().Str("method", "LoadTelegramChatID").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	lg.Trace().Str("username", username).Msgf("%s getting chat ID from Redis...", color.YellowString("➤"))
	chatIDStr, err := rp.rd.Get(ctx, username)
	if err != nil {
		return 0, pkgErrors.Wrapf(err, "get chat ID from Redis by username: %s", username)
	}
	lg.Trace().Str("username", username).Str("chat ID", chatIDStr).Msgf("%s chat ID got from Redis successfully", color.GreenString("✔"))

	lg.Trace().Str("username", username).Str("chat ID", chatIDStr).Msgf("%s converting chat ID to int...", color.YellowString("➤"))
	chatid, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return 0, pkgErrors.Wrapf(err, "convert to int from chat ID: %d", chatid)
	}
	lg.Trace().Str("username", username).Int("chat ID", chatid).Msgf("%s chat ID converted to int successfully", color.GreenString("✔"))

	chatID = int64(chatid)

	return chatID, nil
}
