package rpRedisLoadTelChatID

import (
	"context"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedisLoadTelChatID struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisLoadTelChatID {
	return &RpRedisLoadTelChatID{
		rd: rd,
	}
}

func (rp *RpRedisLoadTelChatID) LoadTelegramChatID(ctx context.Context, username string) (chatID int64, err error) {
	chatIDStr, err := rp.rd.Get(ctx, username)
	if err != nil {
		return 0, pkgErrors.Wrapf(err, "load chat ID from Redis by username: %s", username)
	}

	chatid, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return 0, pkgErrors.Wrapf(err, "convert to int from chat ID: %d", chatid)
	}

	chatID = int64(chatid)

	return chatID, nil
}
