package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/telegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/telegramService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/telegramHandler"
)

type IRepository interface {
	SaveNotice() addNoticeService.IRepository
	DeleteNotice() deleteNoticeService.IRepository
}

type Service struct {
	*addNoticeService.AddNoticeService
	*deleteNoticeService.DeleteNoticeService
	*telegramService.TelegramService
	*consumeNoticeService.ConsumeNoticeService
}

func New(rp IRepository, rb *rabbitmq.Client, tg *telegram.Client) *Service {
	return &Service{
		AddNoticeService:     addNoticeService.New(rp.SaveNotice(), rb),
		DeleteNoticeService:  deleteNoticeService.New(rp.DeleteNotice()),
		TelegramService:      telegramService.New(tg),
		ConsumeNoticeService: consumeNoticeService.New(rb, tg),
	}
}

func (sv *Service) AddNotice() addNoticeHandler.IService {
	return sv.AddNoticeService
}
func (sv *Service) DeleteNotice() deleteNoticeHandler.IService {
	return sv.DeleteNoticeService
}
func (sv *Service) TelegramHandler() telegramHandler.IService {
	return sv.TelegramService
}
