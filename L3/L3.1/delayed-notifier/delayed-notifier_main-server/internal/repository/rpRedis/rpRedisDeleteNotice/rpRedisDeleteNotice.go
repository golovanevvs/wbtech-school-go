package rpRedisDeleteNotice

import (
	"context"

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
	// key := fmt.Sprintf("notices:%d", id)

	// data, err := rp.rd.Get(ctx, key)
	// if err != nil {
	// 	rp.lg.Error().Err(err).Msg("failed get data")
	// 	return nil, err
	// }

	// err = json.Unmarshal([]byte(data), notice)
	// if err != nil {
	// 	rp.lg.Error().Err(err).Msg("failed unmarshal data")
	// 	return nil, err
	// }

	return nil
}
