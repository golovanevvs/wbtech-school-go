package telegramHandler

import (
	"context"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/customerrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	HandleStart(ctx context.Context, username string, chatID int64, message string) error
}

type Handler struct {
	lg zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(rt *ginext.Engine, sv IService) *Handler {
	lg := zlog.Logger.With().Str("component", "handler-telegramHandler").Logger()
	return &Handler{
		lg: lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.POST("/telegram/webhook", hd.WebHookHandler)
}

func (hd *Handler) WebHookHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("handler", "WebHookHandler").Logger()

	lg.Trace().Msg("----- handler is starting")
	defer lg.Trace().Msg("----- handler stopped")

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Msg("invalid content-type")
		c.JSON(http.StatusBadRequest, ginext.H{"error": customerrors.ErrContentTypeAJ.Error()})
		return
	}

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		lg.Warn().Err(err).Msg("error bind json")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	if update.Message == nil {
		c.Status(http.StatusOK)
		return
	}

	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	message := update.Message.Text

	if err := hd.sv.HandleStart(c.Request.Context(), username, chatID, message); err != nil {
		lg.Warn().Int64("chatID", chatID).Str("message", message).Err(err).Msg("failed to handle message")
		c.Status(http.StatusOK)
		return
	}

	lg.Debug().Int64("chatID", chatID).Str("message", message).Msg("Telegram message processed")
	c.Status(http.StatusOK)
}
