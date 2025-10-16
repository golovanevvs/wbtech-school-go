package rpRedisSaveTelName

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedisSaveTelName struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisSaveTelName {
	return &RpRedisSaveTelName{
		rd: rd,
	}
}

func (rp *RpRedisSaveTelName) SaveTelName(ctx context.Context, name string, chatID int64) (err error) {
	err = rp.rd.Set(ctx, name, chatID, 0)
	if err != nil {
		return pkgErrors.Wrapf(err, "save to Redis, name: %s, chat ID: %d", name, chatID)
	}

	return nil
}
