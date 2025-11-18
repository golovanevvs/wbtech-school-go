package service

import (
	"context"
	"image"
	"io"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
)

type iMetaRepository interface {
	CreateImage(ctx context.Context, originalPath, format string) (*model.Image, error)
	CreateImageWithOperations(ctx context.Context, originalPath, format string, options model.ProcessOptions) (*model.Image, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImageStatus(ctx context.Context, id int, status model.ImageStatus) error
	UpdateImageProcessedPath(ctx context.Context, id int, processedPath string) error
	UpdateOriginalPath(ctx context.Context, id int, originalPath string) error
	DeleteImage(ctx context.Context, id int) error
}

type iFileRepository interface {
	SaveOriginalWithID(file io.Reader, id int, originalFilename string) (string, error)
	SaveProcessed(data []byte, originalPath string) (string, error)
	SaveProcessedFromImage(img image.Image, originalPath string) (string, error)
	SaveThumbnail(img image.Image, originalPath string) (string, error)
	DeleteOriginal(path string) error
	DeleteProcessed(path string) error
}

type iQueueProducer interface {
	SendProcessTask(ctx context.Context, imageID string) error
}

type Service struct {
	rpMeta   iMetaRepository
	rpFile   iFileRepository
	producer iQueueProducer
}

func New(rpMeta iMetaRepository, rpFile iFileRepository, producer iQueueProducer) *Service {
	return &Service{
		rpMeta:   rpMeta,
		rpFile:   rpFile,
		producer: producer,
	}
}
