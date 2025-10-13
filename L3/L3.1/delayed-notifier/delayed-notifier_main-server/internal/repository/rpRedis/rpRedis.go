package rpRedis

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
)

type RpRedis struct {
	rd *pkgRedis.Client
}

func New(rd *pkgRedis.Client) *RpRedis {
	return &RpRedis{
		rd: rd,
	}
}
