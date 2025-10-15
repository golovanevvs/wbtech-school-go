package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/publishNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/telegramService"
)

type iRepository interface {
	addNoticeService.IRepository
	deleteNoticeService.IRepository
	consumeNoticeService.IRepository
	telegramService.IRepository
}

type Service struct {
	*addNoticeService.AddNoticeService
	*deleteNoticeService.DeleteNoticeService
	*telegramService.TelegramService
	*consumeNoticeService.ConsumeNoticeService
	*publishNoticeService.PublishNoticeService
}

func New(rp iRepository, rb *pkgRabbitmq.Client, tg *pkgTelegram.Client) *Service {
	return &Service{
		AddNoticeService:     addNoticeService.New(rp, rb),
		DeleteNoticeService:  deleteNoticeService.New(rp),
		TelegramService:      telegramService.New(tg, rp),
		ConsumeNoticeService: consumeNoticeService.New(rb, tg, rp),
		PublishNoticeService: publishNoticeService.New(rb),
	}
}
