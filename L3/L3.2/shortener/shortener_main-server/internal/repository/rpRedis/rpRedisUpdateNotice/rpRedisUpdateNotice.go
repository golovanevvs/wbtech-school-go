package rpRedisUpdateNotice

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgErrors"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRedis"
// 	"github.com/wb-go/wbf/zlog"
// )

// type RpRedisUpdateNotice struct {
// 	lg *zlog.Zerolog
// 	rd *pkgRedis.Client
// }

// func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedisUpdateNotice {
// 	lg := parentLg.With().Str("component", "RpRedisUpdateNotice").Logger()
// 	return &RpRedisUpdateNotice{
// 		lg: &lg,
// 		rd: rd,
// 	}
// }

// func (rp *RpRedisUpdateNotice) UpdateNotice(ctx context.Context, notice *model.Notice) (err error) {
// 	lg := rp.lg.With().Str("method", "UpdateNotice").Logger()
// 	lg.Trace().Msgf("%s method starting", pkgConst.Start)
// 	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

// 	lg.Trace().Msgf("%s marshaling notice...", pkgConst.OpStart)
// 	data, err := json.Marshal(notice)
// 	if err != nil {
// 		return pkgErrors.Wrap(err, "marshal notice")
// 	}
// 	lg.Trace().Msgf("%s notice marshaled successfully", pkgConst.OpSuccess)

// 	key := fmt.Sprintf("notices:%d", notice.ID)

// 	lg.Trace().Str("key", key).Msgf("%s updating notice to Redis...", pkgConst.OpStart)
// 	err = rp.rd.Set(ctx, key, data)
// 	if err != nil {
// 		return pkgErrors.Wrapf(err, "update to Redis, key %s", key)
// 	}
// 	lg.Trace().Str("key", key).Msgf("%s notice updated to Redis successfully", pkgConst.OpSuccess)

// 	return nil
// }
