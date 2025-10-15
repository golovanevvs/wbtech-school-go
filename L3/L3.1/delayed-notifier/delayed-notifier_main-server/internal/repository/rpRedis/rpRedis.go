package rpRedis

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisDeleteNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisLoadNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisLoadTelName"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisSaveNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisSaveTelName"
)

type RpRedis struct {
	*rpRedisSaveNotice.RpRedisSaveNotice
	*rpRedisLoadNotice.RpRedisLoadNotice
	*rpRedisSaveTelName.RpRedisSaveTelName
	*rpRedisLoadTelName.RpRedisLoadTelName
	*rpRedisDeleteNotice.RpRedisDeleteNotice
}

func New(rd *pkgRedis.Client) *RpRedis {
	return &RpRedis{
		RpRedisSaveNotice:   rpRedisSaveNotice.New(rd),
		RpRedisLoadNotice:   rpRedisLoadNotice.New(rd),
		RpRedisSaveTelName:  rpRedisSaveTelName.New(rd),
		RpRedisLoadTelName:  rpRedisLoadTelName.New(rd),
		RpRedisDeleteNotice: rpRedisDeleteNotice.New(rd),
	}
}
