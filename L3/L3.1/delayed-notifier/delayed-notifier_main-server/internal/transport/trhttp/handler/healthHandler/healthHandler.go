package healthHandler

import (
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	lg zlog.Zerolog
	rt *ginext.Engine
}

func New(rt *ginext.Engine) *Handler {
	lg := zlog.Logger.With().Str("component", "healthHandler").Logger()
	return &Handler{
		lg: lg,
		rt: rt,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/healthy", hd.HealthHandler)
}

func (hd *Handler) HealthHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "HealthHandler").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Debug().Msgf("%s healt check successful", pkgConst.OpSuccess)

	c.String(http.StatusOK, "ok")
}
