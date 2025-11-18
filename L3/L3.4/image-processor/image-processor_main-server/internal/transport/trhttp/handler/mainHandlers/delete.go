package mainHandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
)

func (hd *ImageHandlers) DeleteImage(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "DeleteImage").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("param", c.Param("id")).Int("status", http.StatusBadRequest).Msgf("%s invalid image ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, deleteResponse{Error: "invalid image ID"})
		return
	}

	if err := hd.sv.DeleteImage(c.Request.Context(), id); err != nil {
		if errors.Is(err, pkgErrors.ErrNotFound) {
			lg.Warn().Int("image_id", id).Int("status", http.StatusNotFound).Msgf("%s image not found", pkgConst.Warn)
			c.JSON(http.StatusNotFound, deleteResponse{Error: "image not found"})
			return
		}

		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to delete image", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, deleteResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deleteResponse{Message: "image deleted"})

	lg.Debug().Int("image_id", id).Msgf("%s image deleted successfully", pkgConst.OpSuccess)
}
