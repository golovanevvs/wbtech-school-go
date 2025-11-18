package mainHandlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
)

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

func (hd *ImageHandlers) GetAllImages(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetAllImages").Logger()

	images, err := hd.sv.GetAllImages(c.Request.Context())
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to get all images", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]imageResponse, len(images))
	for i, img := range images {
		resp[i] = *convertToImageResponse(&img)
	}

	c.JSON(http.StatusOK, resp)

	lg.Debug().Int("count", len(images)).Msgf("%s all images got successfully", pkgConst.OpSuccess)
}
