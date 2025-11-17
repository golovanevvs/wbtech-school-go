package service

import (
	"context"
	"io"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
)

type iMetaRepository interface {
	CreateImage(ctx context.Context, originalPath, format string) (*model.Image, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImageStatus(ctx context.Context, id int, status model.ImageStatus) error
	UpdateImageProcessedPath(ctx context.Context, id int, processedPath string) error
	DeleteImage(ctx context.Context, id int) error
}

type iFileRepository interface {
	SaveOriginal(file io.Reader, originalFilename string) (string, error)
	SaveProcessed(data []byte, id int) (string, error)
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

func New(repo iMetaRepository, storage iFileRepository, producer iQueueProducer) *Service {
	return &Service{
		rpMeta:   repo,
		rpFile:   storage,
		producer: producer,
	}
}

func (s *Service) UploadImage(ctx context.Context, file io.Reader, filename string) (int, error) {
	format := "jpg"

	originalPath, err := s.rpFile.SaveOriginal(file, filename)
	if err != nil {
		return 0, err
	}

	img, err := s.rpMeta.CreateImage(ctx, originalPath, format)
	if err != nil {
		return 0, err
	}

	err = s.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusProcessing)
	if err != nil {
		return 0, err
	}

	err = s.producer.SendProcessTask(ctx, string(img.ID))
	if err != nil {
		_ = s.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusFailed)
		return 0, err
	}

	return img.ID, nil
}

func (s *Service) GetImage(ctx context.Context, id int) (*model.Image, error) {
	return s.rpMeta.GetImage(ctx, id)
}

func (s *Service) DeleteImage(ctx context.Context, id int) error {
	img, err := s.rpMeta.GetImage(ctx, id)
	if err != nil {
		return err
	}

	_ = s.rpFile.DeleteOriginal(img.OriginalPath)
	_ = s.rpFile.DeleteProcessed(img.ProcessedPath)

	return s.rpMeta.DeleteImage(ctx, id)
}
