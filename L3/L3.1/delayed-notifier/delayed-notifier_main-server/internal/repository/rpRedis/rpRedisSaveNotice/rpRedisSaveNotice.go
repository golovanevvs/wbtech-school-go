package rpRedisSaveNotice

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisSaveNotice struct {
	lg zlog.Zerolog
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisSaveNotice {
	lg := zlog.Logger.With().Str("component", "RpRedisSaveNotice").Logger()
	return &RpRedisSaveNotice{
		lg: lg,
		rd: rd,
	}
}

func (rp *RpRedisSaveNotice) SaveNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	data, err := json.Marshal(notice)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed marshal notice")
		return 0, err
	}

	key, err := rp.rd.SetWithID(ctx, "notices", data, 0)
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed save to Redis")
		return 0, err
	}

	id, err = strconv.Atoi(strings.Split(key, ":")[1])
	if err != nil {
		rp.lg.Error().Err(err).Msg("failed conv")
		return 0, err
	}

	return id, nil
}
