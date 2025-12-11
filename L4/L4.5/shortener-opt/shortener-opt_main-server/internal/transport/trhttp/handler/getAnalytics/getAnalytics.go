package getAnalytics

import (
	"context"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IGetAnalyticsService interface {
	GetAnalytics(ctx context.Context, short string) (totalClicks int, events []model.Analitic, err error)
}

type Handler struct {
	lg             *zlog.Zerolog
	rt             *ginext.Engine
	svGetAnalytics IGetAnalyticsService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, svGetAnalytics IGetAnalyticsService) *Handler {
	lg := parentLg.With().Str("component", "GetAnalytics").Logger()
	return &Handler{
		lg:             &lg,
		rt:             rt,
		svGetAnalytics: svGetAnalytics,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/analytics/:short", hd.GetAnalytics)
}

func (hd *Handler) GetAnalytics(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetAnalytics").Logger()

	short := c.Param("short")
	if short == "" {
		c.JSON(http.StatusBadRequest, response{Error: "short code is required"})
		return
	}

	totalClicks, events, err := hd.svGetAnalytics.GetAnalytics(c.Request.Context(), short)
	if err != nil {
		lg.Error().Err(err).Str("short", short).Msg("Failed to get analytics")
		c.JSON(http.StatusInternalServerError, response{Error: "failed to get analytics"})
		return
	}

	c.JSON(http.StatusOK, response{
		TotalClicks: totalClicks,
		Clicks:      events,
	})

	lg.Debug().Str("short", short).Msgf("%s analutics got successfull", pkgConst.OpSuccess)
}
