package rpFileStorage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type RpFileStorage struct {
	dir string
}

func New(cfg *Config) *RpFileStorage {
	return &RpFileStorage{dir: cfg.Dir}
}

func (fs *RpFileStorage) SaveOriginalWithID(file io.Reader, id int, originalFilename string) (string, error) {
	filename := filepath.Join(fs.dir, fmt.Sprintf("%d_%s", id, originalFilename))
	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}

	return filename, nil
}

func (fs *RpFileStorage) DeleteOriginal(path string) error {
	return os.Remove(path)
}

func (fs *RpFileStorage) DeleteProcessed(path string) error {
	if path != "" {
		return os.Remove(path)
	}
	return nil
}

func (fs *RpFileStorage) SaveProcessedFromImage(img image.Image, originalPath string) (string, error) {
	ext := filepath.Ext(originalPath)
	nameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), ext)

	processedName := strings.Replace(originalPath, nameWithoutExt+ext, nameWithoutExt+"_processed"+ext, 1)

	out, err := os.Create(processedName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(out, img)
	default:
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
	}

	if err != nil {
		return "", err
	}

	return processedName, nil
}

func (fs *RpFileStorage) SaveThumbnail(img image.Image, originalPath string) (string, error) {
	ext := filepath.Ext(originalPath)
	nameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), ext)

	thumbName := strings.Replace(originalPath, nameWithoutExt+ext, nameWithoutExt+"_thumbnail"+ext, 1)

	out, err := os.Create(thumbName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(out, img)
	default:
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
	}

	if err != nil {
		return "", err
	}

	return thumbName, nil
}
