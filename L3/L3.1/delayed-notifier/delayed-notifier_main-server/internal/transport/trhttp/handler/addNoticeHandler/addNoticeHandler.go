package addNoticeHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	AddNotice(ctx context.Context, reqNotice model.ReqNotice) (id int, err error)
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
	hd.rt.POST("/notify", hd.CreateNotice)
}

func (hd *Handler) CreateNotice(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "CreateNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s checking content type...", pkgConst.OpStart)
	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}
	lg.Trace().Msgf("%s content type is valid", pkgConst.OpSuccess)

	lg.Trace().Msgf("%s unmarshaling json data to notice...", pkgConst.OpStart)
	var req model.ReqNotice
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "failed to bind json: " + err.Error()})
		return
	}
	lg.Trace().Msgf("%s json data unmarshaled to notice successfully", pkgConst.OpSuccess)

	lg.Trace().Msgf("%s adding notice...", pkgConst.OpStart)
	id, err := hd.sv.AddNotice(c.Request.Context(), req)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msg("failed to add notice")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to add notice: " + err.Error()})
		return
	}
	lg.Trace().Int("notice ID", id).Msgf("%s notice added successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, ginext.H{"id": id})
}
