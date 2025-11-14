package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/addShortURL"
	"github.com/wb-go/wbf/zlog"
)

type iRepository interface {
	addShortURL.ISaveShortURLRepository
}

type Service struct {
	*addShortURL.Service
}

func New(rs *pkgRetry.Retry, rp iRepository) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()
	return &Service{
		addShortURL.New(&lg, rp),
	}
}
