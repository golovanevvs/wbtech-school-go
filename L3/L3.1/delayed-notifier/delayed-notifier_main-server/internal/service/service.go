package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/consumeNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/getNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/sendNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/telegramStartService"
	"github.com/wb-go/wbf/zlog"
)

type iRepository interface {
	addNoticeService.ISaveNoticeRepository
	deleteNoticeService.IRepository
	getNoticeService.IRepository
	sendNoticeService.IRepository
	telegramStartService.IRepository
}

type Service struct {
	*addNoticeService.AddNoticeService
	*deleteNoticeService.DeleteNoticeService
	*getNoticeService.GetNoticeService
	*telegramStartService.TelegramStartService
	*consumeNoticeService.ConsumeNoticeService
	*sendNoticeService.SendNoticeService
}

func New(rs *pkgRetry.Retry, rp iRepository, rb *pkgRabbitmq.Client, tg *pkgTelegram.Client, em *pkgEmail.Client) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()
	delNotSv := deleteNoticeService.New(rp)
	sendNotSv := sendNoticeService.New(&lg, rs, tg, em, rp)
	return &Service{
		AddNoticeService:     addNoticeService.New(&lg, rb, delNotSv, rp),
		DeleteNoticeService:  delNotSv,
		GetNoticeService:     getNoticeService.New(&lg, rp),
		TelegramStartService: telegramStartService.New(&lg, tg, rp),
		ConsumeNoticeService: consumeNoticeService.New(&lg, rb, delNotSv, sendNotSv),
		SendNoticeService:    sendNotSv,
	}
}
