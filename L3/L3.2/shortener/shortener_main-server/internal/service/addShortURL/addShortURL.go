package addShortURL

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type ISaveShortURLRepository interface {
	SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error)
}

type Service struct {
	lg             *zlog.Zerolog
	rpSaveShortURL ISaveShortURLRepository
}

func New(
	parentLg *zlog.Zerolog,
	rpSaveShortURL ISaveShortURLRepository,
) *Service {
	lg := parentLg.With().Str("component", "AddNoticeService").Logger()
	return &Service{
		lg:             &lg,
		rpSaveShortURL: rpSaveShortURL,
	}
}

func (sv *Service) AddShortURL(ctx context.Context, original, short string) (id int, shortURL string, err error) {
	lg := sv.lg.With().Str("method", "AddShortURL").Logger()
	lg.Debug().Msg("started")

	// createdAt := time.Now()
	// sentAt := reqNotice.SentAt
	// ttl := sentAt.Sub(createdAt)
	// notice := model.Notice{
	// 	UserID:    reqNotice.UserID,
	// 	Message:   reqNotice.Message,
	// 	Channels:  reqNotice.Channels,
	// 	CreatedAt: createdAt,
	// 	SentAt:    sentAt,
	// 	Status:    model.StatusScheduled,
	// }

	// if err := notice.Validate(); err != nil {
	// 	lg.Debug().Err(err).Msgf("%s notice validation failed", pkgConst.Error)
	// 	return 0, pkgErrors.Wrap(err, "notice validation failed")
	// }

	// id, err = sv.rp.SaveNotice(ctx, notice)
	// if err != nil {
	// 	return 0, pkgErrors.Wrap(err, "save notice to repository")
	// }

	// notice.ID = id

	// if err = sv.rb.PublishStructWithTTL(notice, ttl); err != nil {
	// 	lg.Debug().Err(err).Msgf("%s failed to publish struct with TTL to RabbitMQ", pkgConst.Error)
	// 	if err := sv.delNotSv.DeleteNotice(ctx, notice.ID); err != nil {
	// 		lg.Debug().Err(err).Int("notice ID", notice.ID).Msgf("%s failed deleted notice from Redis", pkgConst.Error)
	// 	}
	// 	return 0, pkgErrors.Wrap(err, "error publish struct with TTL to RabbitMQ")
	// }

	return 100, "newShortURL", nil
}
