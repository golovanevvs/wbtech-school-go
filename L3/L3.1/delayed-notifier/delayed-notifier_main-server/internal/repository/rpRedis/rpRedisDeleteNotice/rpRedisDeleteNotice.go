package rpRedisDeleteNotice

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedisDeleteNotice struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisDeleteNotice {
	return &RpRedisDeleteNotice{
		rd: rd,
	}
}

func (rp *RpRedisDeleteNotice) DeleteNotice(ctx context.Context, id int) (err error) {
	key := fmt.Sprintf("notices:%d", id)
	err = rp.rd.Del(ctx, key)
	if err != nil {
		return pkgErrors.Wrapf(err, "delete data from Redis, key: %s", key)
	}

	return nil
}
