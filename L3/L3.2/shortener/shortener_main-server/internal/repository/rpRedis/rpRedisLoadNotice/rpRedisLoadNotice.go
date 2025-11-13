package rpRedisLoadNotice

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgErrors"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRedis"
// 	"github.com/wb-go/wbf/zlog"
// )

// type RpRedisLoadNotice struct {
// 	lg *zlog.Zerolog
// 	rd *pkgRedis.Client
// }

// func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisLoadNotice {
// 	lg := parentLg.With().Str("component", "RpRedisLoadNotice").Logger()
// 	return &RpRedisLoadNotice{
// 		lg: &lg,
// 		rd: rd,
// 	}
// }

// func (rp *RpRedisLoadNotice) LoadNotice(ctx context.Context, id int) (notice *model.Notice, err error) {
// 	lg := rp.lg.With().Str("method", "LoadNotice").Logger()
// 	lg.Trace().Msgf("%s method starting", pkgConst.Start)
// 	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

// 	key := fmt.Sprintf("notices:%d", id)

// 	lg.Trace().Str("key", key).Msgf("%s getting data from Redis...", pkgConst.OpStart)
// 	data, err := rp.rd.Get(ctx, key)
// 	if err != nil {
// 		lg.Trace().Err(err).Int("notice ID", id).Msgf("%s failed to get data from Redis", pkgConst.Error)
// 		return nil, pkgErrors.Wrap(err, "getting data from Redis")
// 	}
// 	lg.Trace().Str("key", key).Msgf("%s data got from Redis successfully", pkgConst.OpSuccess)

// 	lg.Trace().Msgf("%s unmarshaling data to notice...", pkgConst.OpStart)
// 	notice = &model.Notice{}
// 	err = json.Unmarshal([]byte(data), notice)
// 	if err != nil {
// 		lg.Trace().Err(err).Int("notice ID", id).Msgf("%s failed to unmarshal data to notice", pkgConst.Error)
// 		return nil, pkgErrors.Wrap(err, "unmarshal data")
// 	}
// 	lg.Trace().Int("notice ID", notice.ID).Msgf("%s data unmarshaled successfully", pkgConst.OpSuccess)

// 	return notice, nil
// }
