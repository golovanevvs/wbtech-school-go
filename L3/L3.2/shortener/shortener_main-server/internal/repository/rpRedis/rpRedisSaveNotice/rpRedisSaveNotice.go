package rpRedisSaveNotice

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

type RedisClient interface {
	SetWithID(ctx context.Context, prefix string, value interface{}, ttl ...time.Duration) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
}

type RpRedisSaveNotice struct {
	lg *zlog.Zerolog
	rd RedisClient
}

func New(parentLg *zlog.Zerolog, rd RedisClient) *RpRedisSaveNotice {
	lg := parentLg.With().Str("component", "RpRedisSaveNotice").Logger()
	return &RpRedisSaveNotice{
		lg: &lg,
		rd: rd,
	}
}

func (rp *RpRedisSaveNotice) SaveNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	lg := rp.lg.With().Str("method", "SaveNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s marshaling notice...", pkgConst.OpStart)
	data, err := json.Marshal(notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "marshal notice")
	}
	lg.Trace().Msgf("%s notice marshaled successfully", pkgConst.OpSuccess)

	lg.Trace().Msgf("%s saving notice to Redis...", pkgConst.OpStart)
	key, err := rp.rd.SetWithID(ctx, "notices", data)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save to Redis")
	}
	lg.Trace().Str("key", key).Msgf("%s notice saved to Redis successfully", pkgConst.OpSuccess)

	lg.Trace().Msgf("%s converting notice ID to int...", pkgConst.OpStart)
	id, err = strconv.Atoi(strings.Split(key, ":")[1])
	if err != nil {
		return 0, pkgErrors.Wrap(err, "convert string to int")
	}
	lg.Trace().Int("notice ID", id).Msgf("%s notice ID converted to int successfully", pkgConst.OpSuccess)

	notice.ID = id

	lg.Trace().Msgf("%s marshaling notice...", pkgConst.OpStart)
	data, err = json.Marshal(notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "marshal notice")
	}
	lg.Trace().Msgf("%s notice marshaled successfully", pkgConst.OpSuccess)

	lg.Trace().Msgf("%s saving notice to Redis...", pkgConst.OpStart)
	err = rp.rd.Set(ctx, key, data)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save to Redis")
	}
	lg.Trace().Str("key", key).Msgf("%s notice saved to Redis successfully", pkgConst.OpSuccess)

	return id, nil
}
