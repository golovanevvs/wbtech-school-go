package rpRedis_test

// import (
// 	"testing"

// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRedis"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository/rpRedis"
// 	"github.com/stretchr/testify/require"

// 	"github.com/wb-go/wbf/zlog"
// )

// func TestNewRpRedis(t *testing.T) {
// 	mockRedis := &pkgRedis.Client{}

// 	logger := zlog.Logger

// 	r := rpRedis.New(&logger, mockRedis)

// 	require.NotNil(t, r)
// 	require.NotNil(t, r.RpRedisSaveNotice)
// 	require.NotNil(t, r.RpRedisLoadNotice)
// 	require.NotNil(t, r.RpRedisDeleteNotice)
// 	require.NotNil(t, r.RpRedisUpdateNotice)
// 	require.NotNil(t, r.RpRedisSaveChatID)
// 	require.NotNil(t, r.RpRedisLoadTelChatID)
// }
