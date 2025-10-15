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
	rpRedisSaveNotice   *rpRedisSaveNotice.RpRedisSaveNotice
	rpRedisLoadNotice   *rpRedisLoadNotice.RpRedisLoadNotice
	rpRedisSaveTelName  *rpRedisSaveTelName.RpRedisSaveTelName
	rpRedisLoadTelName  *rpRedisLoadTelName.RpRedisLoadTelName
	rpRedisDeleteNotice *rpRedisDeleteNotice.RpRedisDeleteNotice
}

func New(rd *pkgRedis.Client) *RpRedis {
	return &RpRedis{
		rpRedisSaveNotice:   rpRedisSaveNotice.New(rd),
		rpRedisLoadNotice:   rpRedisLoadNotice.New(rd),
		rpRedisSaveTelName:  rpRedisSaveTelName.New(rd),
		rpRedisLoadTelName:  rpRedisLoadTelName.New(rd),
		rpRedisDeleteNotice: rpRedisDeleteNotice.New(rd),
	}
}

func (rp *RpRedis) SaveNotice() *rpRedisSaveNotice.RpRedisSaveNotice {
	return rp.rpRedisSaveNotice
}

func (rp *RpRedis) SaveTelName() *rpRedisSaveTelName.RpRedisSaveTelName {
	return rp.rpRedisSaveTelName
}

func (rp *RpRedis) LoadNotice() *rpRedisLoadNotice.RpRedisLoadNotice {
	return rp.rpRedisLoadNotice
}

func (rp *RpRedis) LoadTelName() *rpRedisLoadTelName.RpRedisLoadTelName {
	return rp.rpRedisLoadTelName
}

func (rp *RpRedis) DeleteNotice() *rpRedisDeleteNotice.RpRedisDeleteNotice {
	return rp.rpRedisDeleteNotice
}
