package rpRedisDeleteNotice

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisDeleteNotice struct {
	lg zlog.Zerolog
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisDeleteNotice {
	lg := zlog.Logger.With().Str("component", "RpRedisDeleteNotice").Logger()
	return &RpRedisDeleteNotice{
		lg: lg,
		rd: rd,
	}
}

func (rp *RpRedisDeleteNotice) DeleteNotice(ctx context.Context, id int) (err error) {
	key := fmt.Sprintf("notices:%d", id)
	err = rp.rd.Del(ctx, key)
	if err != nil {
		rp.lg.Error().Err(err).Str("key", key).Msg("failed deleted data")
		return err
	}

	return nil
}
