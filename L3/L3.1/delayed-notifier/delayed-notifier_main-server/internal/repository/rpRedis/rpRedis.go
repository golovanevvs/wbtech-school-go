package rpRedis

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisDeleteNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisLoadNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisLoadTelChatID"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisSaveNotice"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisSaveTelChatID"
	"github.com/wb-go/wbf/zlog"
)

type RpRedis struct {
	*rpRedisSaveNotice.RpRedisSaveNotice
	*rpRedisLoadNotice.RpRedisLoadNotice
	*rpRedisDeleteNotice.RpRedisDeleteNotice
	*rpRedisSaveTelChatID.RpRedisSaveChatID
	*rpRedisLoadTelChatID.RpRedisLoadTelChatID
}

func New(parentLg *zlog.Zerolog, rd *pkgRedis.Client) *RpRedis {
	lg := parentLg.With().Str("component-1", "RpRedis").Logger()
	return &RpRedis{
		RpRedisSaveNotice:    rpRedisSaveNotice.New(&lg, rd),
		RpRedisLoadNotice:    rpRedisLoadNotice.New(&lg, rd),
		RpRedisDeleteNotice:  rpRedisDeleteNotice.New(&lg, rd),
		RpRedisSaveChatID:    rpRedisSaveTelChatID.New(&lg, rd),
		RpRedisLoadTelChatID: rpRedisLoadTelChatID.New(&lg, rd),
	}
}
