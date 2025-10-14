package rpRedisLoadNotice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisLoadNotice struct {
	lg zlog.Zerolog
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisLoadNotice {
	lg := zlog.Logger.With().Str("component", "RpRedisLoadNotice").Logger()
	return &RpRedisLoadNotice{
		lg: lg,
		rd: rd,
	}
}

func (rp *RpRedisLoadNotice) LoadNotice(ctx context.Context, id int) (notice *model.Notice, err error) {
	key := fmt.Sprintf("notices:%d", id)

	data, err := rp.rd.Get(ctx, key)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed get data")
		return nil, err
	}

	err = json.Unmarshal([]byte(data), notice)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed unmarshal data")
		return nil, err
	}

	return notice, nil
}
