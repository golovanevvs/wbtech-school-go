package mainHandlers

import (
	"context"
	"io"
	"path/filepath"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	UploadImage(ctx context.Context, file io.Reader, filename string, options model.ProcessOptions) (int, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	GetAllImages(ctx context.Context) ([]model.Image, error)
	DeleteImage(ctx context.Context, id int) error
}

type ImageHandlers struct {
	lg             *zlog.Zerolog
	rt             *ginext.Engine
	sv             IService
	fileStorageDir string
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService, fileStorageDir string) *ImageHandlers {
	lg := parentLg.With().Str("component", "ImageProcessor").Logger()
	return &ImageHandlers{
		lg:             &lg,
		rt:             rt,
		sv:             sv,
		fileStorageDir: fileStorageDir,
	}
}

func (hd *ImageHandlers) RegisterRoutes() {
	hd.rt.POST("/upload", hd.UploadImage)
	hd.rt.GET("/image/:id", hd.GetImage)
	hd.rt.GET("/images", hd.GetAllImages)
	hd.rt.DELETE("/image/:id", hd.DeleteImage)

	hd.rt.Static("/uploads", hd.fileStorageDir)
}

func convertToImageResponse(img *model.Image) *imageResponse {
	var processedUrl string
	if img.ProcessedPath != nil {
		processedUrl = "/image-processor_main-server/uploads/" + filepath.Base(*img.ProcessedPath)
	}

	return &imageResponse{
		ID:           img.ID,
		Status:       string(img.Status),
		OriginalPath: img.OriginalPath,
		ProcessedUrl: processedUrl,
		CreatedAt:    img.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Operations:   img.Operations,
	}
}
