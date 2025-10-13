package addNoticeHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/customerrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	AddNotice(ctx context.Context, reqNotice model.ReqNotice) (id int, err error)
}

type Handler struct {
	lg zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(rt *ginext.Engine, sv IService) *Handler {
	lg := zlog.Logger.With().Str("component", "handler-addNoticeHandler").Logger()
	return &Handler{
		lg: lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.POST("/notify", hd.CreateNotice)
}

func (hd *Handler) CreateNotice(c *ginext.Context) {
	lg := hd.lg.With().Str("handler", "createNotice").Logger()

	lg.Trace().Msg("----- handler is starting")
	defer lg.Trace().Msg("----- handler stopped")

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Debug().Str("content-type", c.ContentType()).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": customerrors.ErrContentTypeAJ.Error()})
		return
	}

	var req model.ReqNotice
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Debug().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "failed to bind json: " + err.Error()})
		return
	}

	id, err := hd.sv.AddNotice(c.Request.Context(), req)
	if err != nil {
		lg.Debug().Err(err).Msg("failed to add notice")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "failed to add notice: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id})
}
