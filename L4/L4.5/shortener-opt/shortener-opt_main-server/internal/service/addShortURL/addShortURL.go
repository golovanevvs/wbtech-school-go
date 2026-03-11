package addShortURL

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/model"
	"github.com/google/uuid"
	"github.com/jxskiss/base62"
)

// syncPool for buffer reuse
var byteSlicePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 16)
		return &b
	},
}

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
	// Optimized version: uses rand instead of uuid for fewer allocations
	// and syncPool for buffer reuse

	// Get buffer from pool
	bufPtr := byteSlicePool.Get().(*[]byte)
	buf := *bufPtr
	defer byteSlicePool.Put(bufPtr)

	// Use rand instead of uuid - faster and fewer allocations
	binary.LittleEndian.PutUint64(buf[:8], uint64(time.Now().UnixNano()))
	binary.LittleEndian.PutUint64(buf[8:], uint64(rand.Uint64()))

	return base62.EncodeToString(buf[:])[:8]
}

// generateShortCodeOld - NON-OPTIMIZED version with uuid
// Kept for performance comparison
func (sv *AddShortURLService) generateShortCodeOld() string {
	short := uuid.New()
	return base62.EncodeToString(short[:])[:8]
}
