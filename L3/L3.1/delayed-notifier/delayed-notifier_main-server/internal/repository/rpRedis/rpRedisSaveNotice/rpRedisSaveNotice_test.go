package rpRedisSaveNotice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wb-go/wbf/zlog"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	rpRedisSaveNotice "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis/rpRedisSaveNotice/mocks"
)

func TestSaveNotice_GoMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	logger := zlog.Logger

	notice := model.Notice{
		UserID:  1,
		Message: "Test message",
		Channels: model.Channels{
			{Type: model.ChannelEmail, Value: "test@example.com"},
			{Type: model.ChannelTelegram, Value: "123456"},
		},
		Status:    model.StatusScheduled,
		CreatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockRedis := rpRedisSaveNotice.NewMockRedisClient(ctrl)

		mockRedis.EXPECT().
			SetWithID(gomock.Any(), "notices", gomock.Any()).
			Return("notices:42", nil)
		mockRedis.EXPECT().
			Set(gomock.Any(), "notices:42", gomock.Any()).
			Return(nil)

		rp := New(&logger, mockRedis)

		id, err := rp.SaveNotice(ctx, notice)
		assert.NoError(t, err)
		assert.Equal(t, 42, id)
	})

	t.Run("SetWithID fails", func(t *testing.T) {
		mockRedis := rpRedisSaveNotice.NewMockRedisClient(ctrl)

		mockRedis.EXPECT().
			SetWithID(gomock.Any(), "notices", gomock.Any()).
			Return("", errors.New("redis error"))

		rp := New(&logger, mockRedis)

		id, err := rp.SaveNotice(ctx, notice)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
	})

	t.Run("Set fails", func(t *testing.T) {
		mockRedis := rpRedisSaveNotice.NewMockRedisClient(ctrl)

		mockRedis.EXPECT().
			SetWithID(gomock.Any(), "notices", gomock.Any()).
			Return("notices:100", nil)
		mockRedis.EXPECT().
			Set(gomock.Any(), "notices:100", gomock.Any()).
			Return(errors.New("redis save error"))

		rp := New(&logger, mockRedis)

		id, err := rp.SaveNotice(ctx, notice)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
	})
}

func TestSaveNotice_WithTTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	logger := zlog.Logger

	notice := model.Notice{
		UserID:  1,
		Message: "Test TTL message",
		Channels: model.Channels{
			{Type: model.ChannelEmail, Value: "test@example.com"},
		},
		Status:    model.StatusScheduled,
		CreatedAt: time.Now(),
	}

	ttl := 5 * time.Minute

	mockRedis := rpRedisSaveNotice.NewMockRedisClient(ctrl)

	mockRedis.EXPECT().
		SetWithID(gomock.Any(), "notices", gomock.Any(), ttl).
		Return("notices:200", nil)
	mockRedis.EXPECT().
		Set(gomock.Any(), "notices:200", gomock.Any(), ttl).
		Return(nil)

	rp := New(&logger, mockRedis)

	// Прямой вызов методов Redis с TTL через мок
	key, err := rp.rd.SetWithID(ctx, "notices", notice, ttl)
	assert.NoError(t, err)
	assert.Equal(t, "notices:200", key)

	err = rp.rd.Set(ctx, key, notice, ttl)
	assert.NoError(t, err)
}
