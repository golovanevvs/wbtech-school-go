package getAnalytics

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
)

type ILoadAnalyticsRepository interface {
	LoadAnalytics(ctx context.Context, short string) (totalClicks int, events []model.Analitic, err error)
}

type GetAnalyticsService struct {
	rpGetAnalytics ILoadAnalyticsRepository
}

func New(
	rpGetAnalytics ILoadAnalyticsRepository,
) *GetAnalyticsService {
	return &GetAnalyticsService{
		rpGetAnalytics: rpGetAnalytics,
	}
}

func (sv *GetAnalyticsService) GetAnalytics(ctx context.Context, short string) (totalClicks int, events []model.Analitic, err error) {
	if short == "" {
		err = fmt.Errorf("short code cannot be empty")
		return 0, nil, err
	}

	totalClicks, events, err = sv.rpGetAnalytics.LoadAnalytics(ctx, short)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	return totalClicks, events, nil
}
