package mainHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
)

func (hd *ImageHandlers) UploadImage(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "UploadImage").Logger()

	file, err := c.FormFile("file")
	if err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s no file provided", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, uploadResponse{Error: "no file provided"})
		return
	}

	opsStr := c.PostForm("options")
	var options model.ProcessOptions

	if opsStr != "" {
		if err := json.Unmarshal([]byte(opsStr), &options); err != nil {
			lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s invalid options", pkgConst.Warn)
			c.JSON(http.StatusBadRequest, uploadResponse{Error: "invalid options"})
			return
		}
	} else {
		options = model.ProcessOptions{
			Resize:    true,
			Thumbnail: false,
			Watermark: false,
		}
	}

	src, err := file.Open()
	if err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s could not open file", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, uploadResponse{Error: "could not open file"})
		return
	}
	defer src.Close()

	id, err := hd.sv.UploadImage(c.Request.Context(), src, file.Filename, options)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to upload image", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, uploadResponse{Error: err.Error()})
		return
	}

	resp := uploadResponse{ID: id}
	c.JSON(http.StatusOK, resp)

	lg.Debug().Int("image_id", id).Msgf("%s image uploaded successfully", pkgConst.OpSuccess)
}
