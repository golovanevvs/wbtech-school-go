package healthHandler

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	lg zlog.Zerolog
	rt *ginext.Engine
}

func New(rt *ginext.Engine) *Handler {
	lg := zlog.Logger.With().Str("component", "handler-healthHandler").Logger()
	return &Handler{
		lg: lg,
		rt: rt,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/healthy", hd.HealthHandler)
}

func (hd *Handler) HealthHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("handler", "healthHandler").Logger()

	lg.Trace().Msg("----- handler is starting")
	defer lg.Trace().Msg("----- handler stopped")

	c.String(http.StatusOK, "ok")
}
