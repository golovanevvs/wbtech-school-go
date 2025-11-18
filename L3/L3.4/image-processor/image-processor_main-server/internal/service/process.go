package service

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	"math"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/nfnt/resize"
)

//go:embed watermark.png
var watermarkFile embed.FS

func (sv *Service) ProcessImage(ctx context.Context, imageID string) error {

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

func createThumbnail(img image.Image) image.Image {
	return resize.Resize(150, 150, img, resize.Lanczos3)
}

func addWatermark(img image.Image) image.Image {
	// Загружаем водяной знак
	watermarkFileData, err := watermarkFile.ReadFile("watermark.png")
	if err != nil {
		return img
	}

	watermark, _, err := image.Decode(bytes.NewReader(watermarkFileData))
	if err != nil {
		return img
	}

	// Размеры основного изображения
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Размеры водяного знака
	wmWidth := watermark.Bounds().Dx()
	wmHeight := watermark.Bounds().Dy()

	// Вычисляем коэффициент масштабирования по ширине и высоте
	scaleW := float64(imgWidth) / float64(wmWidth)
	scaleH := float64(imgHeight) / float64(wmHeight)

	// Выбираем минимальный масштаб, чтобы водяной знак полностью помещался
	scale := math.Min(scaleW, scaleH)

	newWmWidth := int(float64(wmWidth) * scale)
	newWmHeight := int(float64(wmHeight) * scale)

	// Масштабируем водяной знак
	scaledWatermark := resize.Resize(uint(newWmWidth), uint(newWmHeight), watermark, resize.Lanczos3)

	// Позиция по центру
	x := (imgWidth - newWmWidth) / 2
	y := (imgHeight - newWmHeight) / 2

	// Накладываем с низкой непрозрачностью (0.2 = 20%)
	result := imaging.Overlay(img, scaledWatermark, image.Point{X: x, Y: y}, 0.2)

	return result
}
