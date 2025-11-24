package service

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgTelegram"
)

// ITelegramStartRp interface for telegram start repository
type ITelegramStartRp interface {
	GetByTelegramUsername(ctx context.Context, username string) (*model.User, error)
	SaveTelegramChatID(ctx context.Context, userID int, chatID *int64) error
}

// ITelegramStartSv interface for telegram start service
type ITelegramStartSv interface {
	Start(ctx context.Context, username string, chatID int64, message string) error
}

// TelegramStartService implements ITelegramStartSv
type TelegramStartService struct {
	tg *pkgTelegram.Client
	rp ITelegramStartRp
}

// NewTelegramStartService creates a new TelegramStartService
func NewTelegramStartService(tg *pkgTelegram.Client, rp ITelegramStartRp) *TelegramStartService {
	return &TelegramStartService{
		tg: tg,
		rp: rp,
	}
}

// Start handles the /start command and saves the chat ID
func (sv *TelegramStartService) Start(ctx context.Context, username string, chatID int64, message string) error {
	err := sv.tg.HandleStart(chatID, message)
	if err != nil {
		return fmt.Errorf("failed to handle start command: %w", err)
	}

	user, err := sv.rp.GetByTelegramUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to find user by telegram username '%s': %w", username, err)
	}

	if user == nil {
		return fmt.Errorf("user with telegram username '%s' not found", username)
	}

	err = sv.rp.SaveTelegramChatID(ctx, user.ID, &chatID)
	if err != nil {
		return fmt.Errorf("failed to save telegram chat ID for user %d: %w", user.ID, err)
	}

	return nil
}
