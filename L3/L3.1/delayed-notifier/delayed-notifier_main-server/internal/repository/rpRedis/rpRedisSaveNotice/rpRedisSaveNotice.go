package rpRedisSaveNotice

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisSaveNotice struct {
	lg *zlog.Zerolog
	rd *pkgRedis.Client
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisSaveNotice {
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
	key, err := rp.rd.SetWithID(ctx, "notices", data, 0)
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

	return id, nil
}
