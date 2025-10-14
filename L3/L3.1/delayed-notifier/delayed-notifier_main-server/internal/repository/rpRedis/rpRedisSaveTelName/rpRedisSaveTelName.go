package rpRedisSaveTelName

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisSaveTelName struct {
	lg zlog.Zerolog
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisSaveTelName {
	lg := zlog.Logger.With().Str("component", "RpRedisSaveTelName").Logger()
	return &RpRedisSaveTelName{
		lg: lg,
		rd: rd,
	}
}

func (rp *RpRedisSaveTelName) SaveTelName(ctx context.Context, name string, chatID int64) (err error) {
	err = rp.rd.Set(ctx, name, chatID, 0)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed save to Redis")
		return err
	}

	return nil
}
