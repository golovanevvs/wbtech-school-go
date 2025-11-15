package getOriginalURL

import (
	"context"
	"fmt"
)

type ILoadOriginalURLRepository interface {
	LoadOriginalURL(ctx context.Context, short string) (original string, err error)
}

type GetOriginalURLService struct {
	rpLoadOriginalURL ILoadOriginalURLRepository
}

func New(
	rpLoadOriginalURL ILoadOriginalURLRepository,
) *GetOriginalURLService {
	return &GetOriginalURLService{
		rpLoadOriginalURL: rpLoadOriginalURL,
	}
}

func (sv *GetOriginalURLService) GetOriginalURL(ctx context.Context, short string) (originalURL string, err error) {
	if short == "" {
		err = fmt.Errorf("short code cannot be empty")
		return "", err
	}

	originalURL, err = sv.rpLoadOriginalURL.LoadOriginalURL(ctx, short)
	if err != nil {
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}

	return originalURL, err
}
