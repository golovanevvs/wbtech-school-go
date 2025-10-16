package rpRedisLoadTelChatID

import (
	"context"
	"fmt"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisLoadTelChatID struct {
	lg zlog.Zerolog
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisLoadTelChatID {
	lg := zlog.Logger.With().Str("component", "RpRedisLoadTelName").Logger()
	return &RpRedisLoadTelChatID{
		lg: lg,
		rd: rd,
	}
}

func (rp *RpRedisLoadTelChatID) LoadTelegramChatID(ctx context.Context, username string) (chatID int64, err error) {
	chatIDStr, err := rp.rd.Get(ctx, username)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed to load chat ID from Redis")
		return 0, fmt.Errorf("failed to load chat ID from Redis: %w", err)
	}

	chatid, err := strconv.Atoi(chatIDStr)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed to convert chat ID to int")
		return 0, fmt.Errorf("failed to convert chat ID to int: %w", err)
	}

	chatID = int64(chatid)

	return chatID, nil
}
