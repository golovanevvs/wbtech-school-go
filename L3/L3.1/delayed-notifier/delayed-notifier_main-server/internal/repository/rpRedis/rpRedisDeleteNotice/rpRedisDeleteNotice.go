package rpRedisDeleteNotice

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisDeleteNotice struct {
	lg *zlog.Zerolog
	rd *pkgRedis.Client
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisDeleteNotice {
	lg := parentLg.With().Str("component1", "RpRedisDeleteNotice").Logger()
	return &RpRedisDeleteNotice{
		lg: &lg,
		rd: rd,
	}
}

func (rp *RpRedisDeleteNotice) DeleteNotice(ctx context.Context, id int) (err error) {
	lg := rp.lg.With().Str("method", "DeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	key := fmt.Sprintf("notices:%d", id)
	lg.Trace().Str("key", key).Msgf("%s deleting notice from Redis...", color.YellowString("➤"))
	err = rp.rd.Del(ctx, key)
	if err != nil {
		return pkgErrors.Wrapf(err, "delete data from Redis, key: %s", key)
	}
	lg.Trace().Str("key", key).Msgf("%s notice deleted from Redis successfully", color.GreenString("✔"))

	return nil
}
