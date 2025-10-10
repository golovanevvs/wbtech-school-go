package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/addNoticeService"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/deleteNoticeService"
)

type IRepository interface {
	SaveNotice() addNoticeService.IRepository
	DeleteNotice() deleteNoticeService.IRepository
}

type IRabbitMQ interface {
	addNoticeService.IRabbitMQ
}

type Service struct {
	*addNoticeService.AddNoticeService
	*deleteNoticeService.DeleteNoticeService
}

func New(rp IRepository, rb IRabbitMQ) *Service {
	return &Service{
		AddNoticeService:    addNoticeService.New(rp.SaveNotice(), rb),
		DeleteNoticeService: deleteNoticeService.New(rp.DeleteNotice()),
	}
}
