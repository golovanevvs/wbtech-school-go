package deleteNoticeHandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	PreDeleteNotice(ctx context.Context, id int) (err error)
}

type Handler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *Handler {
	lg := parentLg.With().Str("component", "deleteNoticeHandler").Logger()

	return &Handler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.DELETE("/notify/:id", hd.DeleteNotice)
}

func (hd *Handler) DeleteNotice(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "DeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	idStr := c.Param("id")
	if idStr == "" {
		lg.Warn().Msgf("%s id param is missing", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		lg.Warn().Str("notice ID", idStr).Msgf("%s id param is not a valid integer", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": "id must be an integer: " + err.Error()})
		return
	}

	lg.Trace().Msgf("%s deleting notice...", pkgConst.OpStart)
	err = hd.sv.PreDeleteNotice(c.Request.Context(), id)
	if errors.Is(err, pkgErrors.ErrNoticeNotFound) {
		lg.Warn().Err(err).Int("notice ID", id).Msgf("%s no exists notice ID", pkgConst.Warn)
		c.JSON(http.StatusNotFound, ginext.H{"error": "notice with ID=" + idStr + " not found: " + err.Error()})
		return
	}
	if err != nil {
		lg.Warn().Err(err).Int("notice ID", id).Msgf("%s failed to delete notice", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to delete notice: " + err.Error()})
		return
	}
	lg.Debug().Int("notice ID", id).Msgf("%s notice deleted successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, ginext.H{"status": "deleted"})
}
