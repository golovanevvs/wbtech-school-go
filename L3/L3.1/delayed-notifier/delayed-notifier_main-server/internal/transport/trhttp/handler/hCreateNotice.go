package handler

import (
	"net/http"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/customerrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type reqNotice struct {
	UserID   int            `json:"user_id" validate:"required"`
	Message  string         `json:"message" validate:"required"`
	Channels model.Channels `json:"channels"`
}

func (hd *Handler) createNotice(c *ginext.Context) {
	lg := zlog.Logger.With().Str("handler", "createNotice").Logger()

	lg.Trace().Msg("handler is starting")

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Debug().Str("content-type", c.ContentType()).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": customerrors.ErrContentTypeAJ.Error()})
		return
	}

	var req reqNotice
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Debug().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "failed to bind json: " + err.Error()})
		return
	}
}
