package rpRedisLoadNotice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedisLoadNotice struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisLoadNotice {
	return &RpRedisLoadNotice{
		rd: rd,
	}
}

func (rp *RpRedisLoadNotice) LoadNotice(ctx context.Context, id int) (notice *model.Notice, err error) {
	key := fmt.Sprintf("notices:%d", id)

	data, err := rp.rd.Get(ctx, key)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "load data from Redis")
	}

	err = json.Unmarshal([]byte(data), notice)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "unmarshal data")
	}

	return notice, nil
}
