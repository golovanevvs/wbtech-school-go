package rpFileStorage

import (
	"fmt"
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
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (fs *RpFileStorage) SaveProcessed(data []byte, originalPath string) (string, error) {
	ext := filepath.Ext(originalPath)
	nameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), ext)

	processedName := strings.Replace(originalPath, nameWithoutExt+ext, nameWithoutExt+"_processed"+ext, 1)

	err := os.WriteFile(processedName, data, 0644)
	if err != nil {
		return "", err
	}

	return processedName, nil
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
