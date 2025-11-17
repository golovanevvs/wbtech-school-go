package service

import (
	"context"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
)

type iMetaRepository interface {
	CreateImage(ctx context.Context, originalPath, format string) (*model.Image, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImageStatus(ctx context.Context, id int, status model.ImageStatus) error
	UpdateImageProcessedPath(ctx context.Context, id int, processedPath string) error
	UpdateOriginalPath(ctx context.Context, id int, originalPath string) error
	DeleteImage(ctx context.Context, id int) error
}

type iFileRepository interface {
	SaveOriginalWithID(file io.Reader, id int, originalFilename string) (string, error)
	SaveProcessed(data []byte, originalPath string) (string, error)
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

func (s *Service) UploadImage(ctx context.Context, file io.Reader, filename string) (int, error) {
	format := getFileFormat(filename)

	img, err := s.rpMeta.CreateImage(ctx, "", format)
	if err != nil {
		return 0, err
	}

	originalPath, err := s.rpFile.SaveOriginalWithID(file, img.ID, filename)
	if err != nil {
		_ = s.rpMeta.DeleteImage(ctx, img.ID)
		return 0, err
	}

	err = s.rpMeta.UpdateOriginalPath(ctx, img.ID, originalPath)
	if err != nil {
		_ = s.rpFile.DeleteOriginal(originalPath)
		_ = s.rpMeta.DeleteImage(ctx, img.ID)
		return 0, err
	}

	err = s.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusProcessing)
	if err != nil {
		return 0, err
	}

	err = s.producer.SendProcessTask(ctx, strconv.Itoa(img.ID))
	if err != nil {
		_ = s.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusFailed)
		return 0, err
	}

	return img.ID, nil
}

func getFileFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	default:
		return "unknown"
	}
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

	if img.ProcessedPath != nil {
		_ = s.rpFile.DeleteProcessed(*img.ProcessedPath)
	}

	return s.rpMeta.DeleteImage(ctx, id)
}
