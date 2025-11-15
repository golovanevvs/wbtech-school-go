package getOriginalURL

import (
	"context"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IGetOriginalURLService interface {
	GetOriginalURL(ctx context.Context, short string) (originalURL string, err error)
}

type IAddClickEventService interface {
	AddClickEvent(ctx context.Context, clickEvent model.Analitic) (err error)
}

type Handler struct {
	lg               *zlog.Zerolog
	rt               *ginext.Engine
	svGetOriginalURL IGetOriginalURLService
	svAddClickEvent  IAddClickEventService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, svGetShortURL IGetOriginalURLService, svAddClickEvent IAddClickEventService) *Handler {
	lg := parentLg.With().Str("component", "GetOriginalURL").Logger()
	return &Handler{
		lg:               &lg,
		rt:               rt,
		svGetOriginalURL: svGetShortURL,
		svAddClickEvent:  svAddClickEvent,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.GET("/s/:short", hd.RedirectToOriginal)
}

func (hd *Handler) RedirectToOriginal(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetOriginalURL").Logger()

	short := c.Param("short")
	if short == "" {
		c.JSON(http.StatusBadRequest, response{Error: "short code is required"})
		return
	}

	originalURL, err := hd.svGetOriginalURL.GetOriginalURL(c.Request.Context(), short)
	if err != nil {
		lg.Error().Err(err).Str("short", short).Msg("Failed to get original URL")
		c.JSON(http.StatusNotFound, response{Error: "URL not found"})
		return
	}

	clickEvent := model.Analitic{
		Short:     short,
		UserAgent: c.GetHeader("User-Agent"),
		IP:        c.ClientIP(),
	}

	go func() {
		ctx := context.Background()
		if err := hd.svAddClickEvent.AddClickEvent(ctx, clickEvent); err != nil {
			lg.Error().Err(err).Str("short", short).Msg("Failed to save click event asynchronously")
		}
	}()

	c.Redirect(http.StatusFound, originalURL)

	lg.Debug().Str("short", short).Str("originalURL", originalURL).Msgf("%s redirect successfull", pkgConst.OpSuccess)
}
