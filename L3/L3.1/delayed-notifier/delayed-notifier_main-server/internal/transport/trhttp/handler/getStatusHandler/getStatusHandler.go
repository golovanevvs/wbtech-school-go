package getStatusHandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	GetNotice(ctx context.Context, id int) (notice *model.Notice, err error)
}

type Handler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *Handler {
	lg := parentLg.With().Str("component", "addNoticeHandler").Logger()
	return &Handler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/notify/:id", hd.GetNotice)
}

func (hd *Handler) GetNotice(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetNotice").Logger()
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

	lg.Trace().Msgf("%s getting notice...", pkgConst.OpStart)
	notice, err := hd.sv.GetNotice(c.Request.Context(), id)
	if errors.Is(err, pkgErrors.ErrNoticeNotFound) {
		lg.Warn().Err(err).Int("notice ID", id).Msgf("%s notice ID", pkgConst.Warn)
		c.JSON(http.StatusNotFound, ginext.H{"error": "notice with ID=" + idStr + " not found: " + err.Error()})
		return
	}
	if err != nil {
		lg.Warn().Err(err).Int("notice ID", id).Msgf("%s failed to get notice", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to get notice: " + err.Error()})
		return
	}
	lg.Trace().Int("notice ID", id).Msgf("%s notice got successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, ginext.H{"status": notice.Status})
}
