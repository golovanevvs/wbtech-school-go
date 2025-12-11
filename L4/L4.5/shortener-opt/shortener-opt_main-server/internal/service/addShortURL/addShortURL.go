package addShortURL

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/model"
	"github.com/google/uuid"
	"github.com/jxskiss/base62"
)

type ISaveShortURLRepository interface {
	SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error)
}

type AddShortURLService struct {
	rpSaveShortURL ISaveShortURLRepository
}

func New(
	rpSaveShortURL ISaveShortURLRepository,
) *AddShortURLService {
	return &AddShortURLService{
		rpSaveShortURL: rpSaveShortURL,
	}
}

func (sv *AddShortURLService) AddShortURL(ctx context.Context, original, short string) (id int, shortURL string, err error) {
	custom := short != ""
	if !custom {
		short = sv.generateShortCode()
	}

	if original == "" {
		err = fmt.Errorf("original URL cannot be empty")
		return 0, "", err
	}

	shortURLModel := model.ShortURL{
		Original:  original,
		Short:     short,
		Custom:    custom,
		CreatedAt: time.Time{},
	}

	id, err = sv.rpSaveShortURL.SaveShortURL(ctx, shortURLModel)
	if err != nil {
		return 0, "", fmt.Errorf("failed to save short URL: %w", err)
	}

	return id, short, nil
}

func (sv *AddShortURLService) generateShortCode() string {
	short := uuid.New()
	return base62.EncodeToString(short[:])[:8]
}
