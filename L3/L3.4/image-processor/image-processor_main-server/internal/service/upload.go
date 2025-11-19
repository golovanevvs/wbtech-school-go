package service

import (
	"context"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
)

func (sv *Service) UploadImage(ctx context.Context, file io.Reader, filename string, options model.ProcessOptions) (int, error) {
	format := getFileFormat(filename)

	img, err := sv.rpMeta.CreateImage(ctx, "", format, options)
	if err != nil {
		return 0, err
	}

	originalPath, err := sv.rpFile.SaveOriginalWithID(file, img.ID, filename)
	if err != nil {
		_ = sv.rpMeta.DeleteImage(ctx, img.ID)
		return 0, err
	}

	err = sv.rpMeta.UpdateOriginalPath(ctx, img.ID, originalPath)
	if err != nil {
		_ = sv.rpFile.DeleteOriginal(originalPath)
		_ = sv.rpMeta.DeleteImage(ctx, img.ID)
		return 0, err
	}

	err = sv.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusQueued)
	if err != nil {
		return 0, err
	}

	err = sv.producer.SendProcessTask(ctx, strconv.Itoa(img.ID))
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, img.ID, model.StatusFailed)
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
