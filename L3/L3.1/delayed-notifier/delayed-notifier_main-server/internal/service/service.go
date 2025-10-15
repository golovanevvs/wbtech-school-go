package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/telegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/rpRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/telegramService"
)

type iRepository interface {
	RpRedis() *rpRedis.RpRedis
}

type Service struct {
	AddNoticeService     *addNoticeService.AddNoticeService
	DeleteNoticeService  *deleteNoticeService.DeleteNoticeService
	TelegramService      *telegramService.TelegramService
	ConsumeNoticeService *consumeNoticeService.ConsumeNoticeService
}

func New(rp iRepository, rb *rabbitmq.Client, tg *telegram.Client) *Service {
	return &Service{
		AddNoticeService:     addNoticeService.New(rp.RpRedis().SaveNotice(), rb),
		DeleteNoticeService:  deleteNoticeService.New(rp.RpRedis().DeleteNotice()),
		TelegramService:      telegramService.New(tg, rp.RpRedis().SaveTelName()),
		ConsumeNoticeService: consumeNoticeService.New(rb, tg, rp.RpRedis().LoadTelName()),
	}
}
