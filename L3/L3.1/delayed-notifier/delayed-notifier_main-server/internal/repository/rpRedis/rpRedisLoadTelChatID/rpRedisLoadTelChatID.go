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

func (rp *RpRedisLoadTelChatID) LoadTelegramChatID(ctx context.Context, username string) (chatID int, err error) {
	chatIDStr, err := rp.rd.Get(ctx, username)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed to save to Redis")
		return 0, fmt.Errorf("failed to save to Redis: %w", err)
	}

	chatID, err = strconv.Atoi(chatIDStr)

	return chatID, nil
}
