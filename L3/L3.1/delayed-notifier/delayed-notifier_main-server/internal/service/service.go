package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository/postgres"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/telegramService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/addNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/deleteNoticeHandler"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport/trhttp/handler/telegramHandler"
)

type iRepository interface {
	Postgres() *postgres.Postgres
}

type Service struct {
	AddNoticeService     *addNoticeService.AddNoticeService
	DeleteNoticeService  *deleteNoticeService.DeleteNoticeService
	TelegramService      *telegramService.TelegramService
	ConsumeNoticeService *consumeNoticeService.ConsumeNoticeService
}

func New(rp iRepository, rb *pkgRabbitmq.Client, tg *pkgTelegram.Client, rd *pkgRedis.Client) *Service {
	return &Service{
		AddNoticeService:     addNoticeService.New(rp.Postgres().SaveNoticePostgres, rb),
		DeleteNoticeService:  deleteNoticeService.New(rp.Postgres().DeleteNoticePostgres),
		TelegramService:      telegramService.New(tg, rd),
		ConsumeNoticeService: consumeNoticeService.New(rb, tg, rd),
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
