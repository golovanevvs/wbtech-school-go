package healthHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgPrometheus"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine) *Handler {
	lg := parentLg.With().Str("component", "healthHandler").Logger()
	return &Handler{
		lg: &lg,
		rt: rt,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/healthy", hd.HealthHandler)
	hd.rt.GET("/metrics", gin.WrapH(pkgPrometheus.Handler()))
}

func (hd *Handler) HealthHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "HealthHandler").Logger()

	lg.Debug().Msgf("%s healt check successful", pkgConst.OpSuccess)

	c.String(http.StatusOK, "ok")
}
