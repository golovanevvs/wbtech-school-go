package telegramHandler

import (
	"context"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
	Start(ctx context.Context, username string, chatID int64, message string) error
}

type Handler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *Handler {
	lg := parentLg.With().Str("component", "handler-telegramHandler").Logger()
	return &Handler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *Handler) RegisterRoutes() {
	hd.rt.POST("/telegram/webhook", hd.WebHookHandler)
}

func (hd *Handler) WebHookHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("handler", "WebHookHandler").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Msgf("%s checking content type...", pkgConst.OpStart)
	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}
	lg.Trace().Msgf("%s content type is valid", pkgConst.OpSuccess)

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
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

	if err := hd.sv.Start(c.Request.Context(), username, chatID, message); err != nil {
		lg.Warn().Int64("chatID", chatID).Str("message_body", message).Err(err).Msgf("%s failed to handle message", pkgConst.Warn)
		c.Status(http.StatusOK)
		return
	}

	lg.Debug().Int64("chatID", chatID).Str("message_body", message).Msgf("%s Telegram message processed", pkgConst.OpSuccess)
	c.Status(http.StatusOK)
}
