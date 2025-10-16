package rpRedisSaveNotice

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedisSaveNotice struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedisSaveNotice {
	return &RpRedisSaveNotice{
		rd: rd,
	}
}

func (rp *RpRedisSaveNotice) SaveNotice(ctx context.Context, notice model.Notice) (id int, err error) {
	data, err := json.Marshal(notice)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "marshal notice")
	}

	key, err := rp.rd.SetWithID(ctx, "notices", data, 0)
	if err != nil {
		return 0, pkgErrors.Wrap(err, "save to Redis")
	}

	id, err = strconv.Atoi(strings.Split(key, ":")[1])
	if err != nil {
		return 0, pkgErrors.Wrap(err, "convert string to int")
	}

	return id, nil
}
