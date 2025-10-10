package deleteNoticeHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/customerrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	DeleteNotice(ctx context.Context, id int) (err error)
}

type Handler struct {
	rt *ginext.Engine
	sv IService
}

func New(rt *ginext.Engine, sv IService) *Handler {
	return &Handler{
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.DELETE("/notify", hd.DeleteNotice)
}

func (hd *Handler) DeleteNotice(c *ginext.Context) {
	lg := zlog.Logger.With().Str("handler", "deleteNotice").Logger()

	lg.Trace().Msg("handler is starting")

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Debug().Str("content-type", c.ContentType()).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": customerrors.ErrContentTypeAJ.Error()})
		return
	}
}
