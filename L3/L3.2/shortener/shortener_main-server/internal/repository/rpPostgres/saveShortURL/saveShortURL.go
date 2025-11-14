package saveShortURL

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type RpPostgresSaveShortURL struct {
	lg *zlog.Zerolog
}

func New(parentLg *zlog.Zerolog) *RpPostgresSaveShortURL {
	lg := parentLg.With().Str("component", "RpRedisDeleteNotice").Logger()
	return &RpPostgresSaveShortURL{
		lg: &lg,
	}
}

func (rp *RpPostgresSaveShortURL) SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error) {
	return 100, nil
}
