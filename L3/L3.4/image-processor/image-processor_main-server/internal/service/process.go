package service

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
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

	processedImg := sv.applyProcessing(imgDecoded, img.Operations)

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

func (sv *Service) applyProcessing(img image.Image, options model.ProcessOptions) image.Image {
	result := img

	if options.Resize {
		result = applyClassicCartoonEffect(result)
	}
	// if options.Resize {
	// 	result = imaging.Fit(result, 800, 600, imaging.Lanczos)
	// 	result = imaging.Fill(result, 800, 600, imaging.Center, imaging.Lanczos)
	// }

	if options.Watermark {
		result = addWatermark(result)
	}

	if options.Thumbnail {
		result = createThumbnail(result)
	}

	return result
}

func decodeImage(file *os.File) (image.Image, string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func createThumbnail(img image.Image) image.Image {
	scaled := imaging.Fit(img, 150, 150, imaging.Lanczos)

	background := imaging.New(150, 150, color.RGBA{255, 255, 255, 255})

	result := imaging.PasteCenter(background, scaled)

	return result
}

func addWatermark(img image.Image) image.Image {
	watermarkFileData, err := watermarkFile.ReadFile("watermark.png")
	if err != nil {
		return img
	}

	watermark, _, err := image.Decode(bytes.NewReader(watermarkFileData))
	if err != nil {
		return img
	}

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	wmWidth := watermark.Bounds().Dx()
	wmHeight := watermark.Bounds().Dy()

	scaleW := float64(imgWidth) / float64(wmWidth)
	scaleH := float64(imgHeight) / float64(wmHeight)

	scale := math.Max(scaleW, scaleH)

	newWmWidth := int(float64(wmWidth) * scale)
	newWmHeight := int(float64(wmHeight) * scale)

	scaledWatermark := resize.Resize(uint(newWmWidth), uint(newWmHeight), watermark, resize.Lanczos3)

	x := (imgWidth - newWmWidth) / 2
	y := (imgHeight - newWmHeight) / 2

	result := imaging.Overlay(img, scaledWatermark, image.Point{X: x, Y: y}, 0.2)

	return result
}

// applyClassicCartoonEffect применяет классические фильтры для мультяшного эффекта
func applyClassicCartoonEffect(img image.Image) image.Image {
	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return img
	}
	defer mat.Close()

	// 1. Упрощение деталей с помощью bilateral filter
	simplified := gocv.NewMat()
	defer simplified.Close()
	gocv.BilateralFilter(mat, &simplified, 9, 150, 150)

	// 2. Дополнительное упрощение median filter
	smoothed := gocv.NewMat()
	defer smoothed.Close()
	gocv.MedianBlur(simplified, &smoothed, 7)

	// 3. Создаем маску краев
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(smoothed, &gray, gocv.ColorBGRToGray)

	edges := gocv.NewMat()
	defer edges.Close()
	gocv.AdaptiveThreshold(gray, &edges, 255, gocv.AdaptiveThresholdGaussian,
		gocv.ThresholdBinary, 9, 2)

	// 4. Инвертируем края для создания маски
	edgesInv := gocv.NewMat()
	defer edgesInv.Close()
	gocv.BitwiseNot(edges, &edgesInv)

	// 5. Создаем цветную версию инвертированных краев
	edgesInvColor := gocv.NewMat()
	defer edgesInvColor.Close()
	gocv.CvtColor(edgesInv, &edgesInvColor, gocv.ColorGrayToBGR)

	// 6. Объединяем упрощенное изображение с маской
	result := gocv.NewMat()
	defer result.Close()
	gocv.BitwiseAnd(smoothed, edgesInvColor, &result)

	resultImg, err := result.ToImage()
	if err != nil {
		return img
	}

	return resultImg
}
