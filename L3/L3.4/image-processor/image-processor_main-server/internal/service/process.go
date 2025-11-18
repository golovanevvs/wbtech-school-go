package service

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/nfnt/resize"
)

func (sv *Service) ProcessImage(ctx context.Context, imageID string) error {
	log.Printf("Processing image: %s", imageID)

	id, err := strconv.Atoi(imageID)
	if err != nil {
		return fmt.Errorf("invalid image ID: %w", err)
	}

	img, err := sv.rpMeta.GetImage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	err = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusProcessing)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	originalFile, err := os.Open(img.OriginalPath)
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusFailed)
		return fmt.Errorf("failed to open original file: %w", err)
	}
	defer originalFile.Close()

	imgDecoded, _, err := decodeImage(originalFile)
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusFailed)
		return fmt.Errorf("failed to decode image: %w", err)
	}

	processedImg, err := sv.applyProcessingWithOptions(imgDecoded, img.Operations, img.OriginalPath)
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusFailed)
		return fmt.Errorf("failed to apply processing: %w", err)
	}

	processedPath, err := sv.rpFile.SaveProcessedFromImage(processedImg, img.OriginalPath)
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusFailed)
		return fmt.Errorf("failed to save processed file: %w", err)
	}

	err = sv.rpMeta.UpdateImageProcessedPath(ctx, id, processedPath)
	if err != nil {
		_ = sv.rpMeta.UpdateImageStatus(ctx, id, model.StatusFailed)
		return fmt.Errorf("failed to update processed path: %w", err)
	}

	return nil
}

func (sv *Service) applyProcessingWithOptions(img image.Image, options model.ProcessOptions, originalPath string) (image.Image, error) {
	result := img

	if options.Resize {
		result = resize.Resize(800, 600, result, resize.Lanczos3)
	}

	if options.Watermark {
		result = addWatermark(result)
	}

	if options.Thumbnail {
		thumbnail := createThumbnail(result)
		_, err := sv.rpFile.SaveThumbnail(thumbnail, originalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to save thumbnail: %w", err)
		}
	}

	return result, nil
}

func decodeImage(file *os.File) (image.Image, string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func addWatermark(img image.Image) image.Image {
	bounds := img.Bounds()
	wmImg := image.NewRGBA(bounds)
	draw.Draw(wmImg, bounds, img, bounds.Min, draw.Src)

	// Простой водяной знак — прямоугольник в правом нижнем углу
	wmBounds := image.Rect(bounds.Dx()-100, bounds.Dy()-30, bounds.Dx(), bounds.Dy())
	draw.Draw(wmImg, wmBounds, &image.Uniform{color.RGBA{255, 255, 255, 128}}, image.Point{}, draw.Over)

	return wmImg
}

func createThumbnail(img image.Image) image.Image {
	return resize.Resize(150, 150, img, resize.Lanczos3)
}
