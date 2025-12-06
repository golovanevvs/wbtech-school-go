package telegramHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// ISvForTelegramHandler interface for telegram handler service
type ISvForTelegramHandler interface {
	Start(ctx context.Context, username string, chatID int64, message string) error
}

// TelegramHandler handles Telegram webhook requests
type TelegramHandler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv ISvForTelegramHandler
}

// NewTelegramHandler creates a new TelegramHandler
func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv ISvForTelegramHandler) *TelegramHandler {
	lg := parentLg.With().Str("component", "handler-telegramHandler").Logger()
	return &TelegramHandler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

// RegisterRoutes registers the Telegram webhook route
func (hd *TelegramHandler) RegisterRoutes() {
	hd.rt.POST("/telegram/webhook", hd.WebHookHandler)
}

// WebHookHandler handles incoming Telegram webhook requests
func (hd *TelegramHandler) WebHookHandler(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "WebHookHandler").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if update.Message == nil {
		c.Status(http.StatusOK)
		return
	}

	username := update.Message.From.UserName
	if username == "" {
		lg.Warn().Int64("chatID", update.Message.Chat.ID).Msg("Username is empty, cannot process /start command")
		c.Status(http.StatusOK)
		return
	}

	chatID := update.Message.Chat.ID
	message := update.Message.Text

	if update.Message.IsCommand() && update.Message.Command() == "start" {
		lg.Debug().Str("username", username).Int64("chatID", chatID).Msg("Processing /start command")

		if err := hd.sv.Start(c.Request.Context(), username, chatID, message); err != nil {
			lg.Warn().Str("username", username).Int64("chatID", chatID).Err(err).Msgf("%s failed to handle /start command", pkgConst.Warn)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lg.Debug().Str("username", username).Int64("chatID", chatID).Msgf("%s Telegram /start command processed successfully", pkgConst.OpSuccess)
		c.Status(http.StatusOK)
		return
	}

	lg.Debug().Str("username", username).Int64("chatID", chatID).Str("message_body", message).Msgf("%s Telegram message processed", pkgConst.OpSuccess)
	c.Status(http.StatusOK)
}
