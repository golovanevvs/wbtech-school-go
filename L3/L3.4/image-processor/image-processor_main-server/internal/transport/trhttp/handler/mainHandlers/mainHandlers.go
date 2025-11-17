package mainHandlers

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (int, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	DeleteImage(ctx context.Context, id int) error
}

type ImageHandlers struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *ImageHandlers {
	lg := parentLg.With().Str("component", "ImageProcessor").Logger()
	return &ImageHandlers{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *ImageHandlers) RegisterRoutes() {
	hd.rt.POST("/upload", hd.UploadImage)
	hd.rt.GET("/image/:id", hd.GetImage)
	hd.rt.DELETE("/image/:id", hd.DeleteImage)
}

func (hd *ImageHandlers) UploadImage(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "UploadImage").Logger()

	file, err := c.FormFile("file")
	if err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s no file provided", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, uploadResponse{Error: "no file provided"})
		return
	}

	src, err := file.Open()
	if err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s could not open file", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, uploadResponse{Error: "could not open file"})
		return
	}
	defer src.Close()

	id, err := hd.sv.UploadImage(c.Request.Context(), src, file.Filename)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to upload image", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, uploadResponse{Error: err.Error()})
		return
	}

	resp := uploadResponse{ID: id}
	c.JSON(http.StatusOK, resp)

	lg.Debug().Int("image_id", id).Msgf("%s image uploaded successfully", pkgConst.OpSuccess)
}

func (hd *ImageHandlers) GetImage(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetImage").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("param", c.Param("id")).Int("status", http.StatusBadRequest).Msgf("%s invalid image ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, imageResponse{Error: "invalid image ID"})
		return
	}

	img, err := hd.sv.GetImage(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to get image", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, imageResponse{Error: err.Error()})
		return
	}

	resp := convertToImageResponse(img)
	c.JSON(http.StatusOK, resp)

	lg.Debug().Int("image_id", id).Msgf("%s image got successfully", pkgConst.OpSuccess)
}

func (hd *ImageHandlers) DeleteImage(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "DeleteImage").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("param", c.Param("id")).Int("status", http.StatusBadRequest).Msgf("%s invalid image ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, deleteResponse{Error: "invalid image ID"})
		return
	}

	if err := hd.sv.DeleteImage(c.Request.Context(), id); err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to delete image", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, deleteResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deleteResponse{Message: "image deleted"})

	lg.Debug().Int("image_id", id).Msgf("%s image deleted successfully", pkgConst.OpSuccess)
}

func convertToImageResponse(img *model.Image) *imageResponse {
	var processedPath string
	if img.ProcessedPath != nil {
		processedPath = *img.ProcessedPath
	}

	return &imageResponse{
		ID:            img.ID,
		Status:        string(img.Status),
		OriginalPath:  img.OriginalPath,
		ProcessedPath: processedPath,
		CreatedAt:     img.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
