package deleteNoticeHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type Handler struct {
	lg zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(rt *ginext.Engine, sv IService) *Handler {
	lg := zlog.Logger.With().Str("component", "deleteNoticeHandler").Logger()

	return &Handler{
		lg: lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.DELETE("/notify", hd.DeleteNotice)
}

func (hd *Handler) DeleteNotice(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "DeleteNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s checking content type...", pkgConst.OpStart)
	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}
	lg.Trace().Msgf("%s content type is valid", pkgConst.OpSuccess)
}
