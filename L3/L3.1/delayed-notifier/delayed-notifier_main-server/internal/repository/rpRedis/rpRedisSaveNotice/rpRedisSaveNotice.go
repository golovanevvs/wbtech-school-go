package rpRedisSaveNotice

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/wb-go/wbf/zlog"
)

type RpRedisSaveNotice struct {
	lg *zlog.Zerolog
	rd *pkgRedis.Client
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisSaveNotice {
	lg := parentLg.With().Str("component-2", "RpRedisSaveNotice").Logger()
	return &RpRedisSaveNotice{
		lg: &lg,
		rd: rd,
	}
}

func (rp *RpRedisSaveNotice) SaveNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	lg := rp.lg.With().Str("method", "SaveNotice").Logger()
	lg.Trace().Msgf("%s method starting", color.GreenString("🟢"))
	defer lg.Trace().Msgf("%s method stopped", color.RedString("🟢"))

	lg.Trace().Msgf("%s marshaling notice...", color.YellowString("➤"))
	data, err := json.Marshal(notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "marshal notice")
	}
	lg.Trace().Msgf("%s notice marshaled successfully", color.GreenString("✔"))

	lg.Trace().Msgf("%s saving notice to Redis...", color.YellowString("➤"))
	key, err := rp.rd.SetWithID(ctx, "notices", data, 0)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save to Redis")
	}
	lg.Trace().Str("key", key).Msgf("%s notice saved to Redis successfully", color.GreenString("✔"))

	lg.Trace().Msgf("%s converting notice ID to int...", color.YellowString("➤"))
	id, err = strconv.Atoi(strings.Split(key, ":")[1])
	if err != nil {
		return 0, pkgErrors.Wrap(err, "convert string to int")
	}
	lg.Trace().Int("notice ID", id).Msgf("%s notice ID converted to int successfully", color.GreenString("✔"))

	return id, nil
}
