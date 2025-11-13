package addShortURL

import (
	"context"
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IAddShortURLService interface {
	AddShortURL(ctx context.Context, original, short string) (id int, shortURL string, err error)
}

type Handler struct {
	lg            *zlog.Zerolog
	rt            *ginext.Engine
	svAddShortURL IAddShortURLService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, svAddShortURL IAddShortURLService) *Handler {
	lg := parentLg.With().Str("component", "AddShortURL").Logger()
	return &Handler{
		lg:            &lg,
		rt:            rt,
		svAddShortURL: svAddShortURL,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.POST("/shorten", hd.AddShortURL)
}

func (hd *Handler) AddShortURL(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "AddShortURL").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Err(pkgErrors.ErrContentTypeAJ).Str("Content-Type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Int("status", http.StatusBadRequest).Msgf("%s failed to bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, response{Error: pkgErrors.ErrBindJSON.Error()})
		return
	}

	id, shortURL, err := hd.svAddShortURL.AddShortURL(c.Request.Context(), req.Original, req.Short)
	if err != nil {
		lg.Error().Err(err).Int("status", http.StatusInternalServerError).Msgf("%s failed to add short URL", pkgConst.Error)
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}
	lg.Debug().Int("ID", id).Str("original", req.Original).Str("shortURL", shortURL).Msgf("%s short URL added successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, response{ShortURL: shortURL})
}
