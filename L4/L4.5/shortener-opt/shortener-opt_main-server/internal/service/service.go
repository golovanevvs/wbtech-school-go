package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/service/addClickEvent"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/service/addShortURL"
	getAnalytics "github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/service/getAnalitycs"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/service/getOriginalURL"
)

type iRepository interface {
	addShortURL.ISaveShortURLRepository
	getOriginalURL.ILoadOriginalURLRepository
	getAnalytics.ILoadAnalyticsRepository
	addClickEvent.ISaveClickEventRepository
}

type Service struct {
	*addShortURL.AddShortURLService
	*getOriginalURL.GetOriginalURLService
	*addClickEvent.AddClickEventService
	*getAnalytics.GetAnalyticsService
}

func New(rs *pkgRetry.Retry, rp iRepository) *Service {
	return &Service{
		AddShortURLService:    addShortURL.New(rp),
		GetOriginalURLService: getOriginalURL.New(rp),
		GetAnalyticsService:   getAnalytics.New(rp),
		AddClickEventService:  addClickEvent.New(rp),
	}
}
