package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/transport/trhttp/handler/telegramHandler"
	"github.com/wb-go/wbf/zlog"
)

// Service structure that combines all services
type Service struct {
	TelegramStart *TelegramStartService
	Notice        *NoticeService
}

// New creates a new Service structure
func New(
	cfg *Config,
	rp *repository.Repository,
	tgClient *pkgTelegram.Client,
	emailClient *pkgEmail.Client,
	rs *pkgRetry.Retry,
) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()

	telegramStartService := NewTelegramStartService(tgClient, rp.User())
	noticeService := NewNoticeService(&lg, rs, tgClient, emailClient, rp.NoticeRepository)

	return &Service{
		TelegramStart: telegramStartService,
		Notice:        noticeService,
	}
}

// TelegramService returns the telegram service
func (sv *Service) TelegramService() telegramHandler.ISvForTelegramHandler {
	return sv.TelegramStart
}
