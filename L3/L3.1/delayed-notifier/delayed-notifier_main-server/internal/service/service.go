package service

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service/noticeservice"
	"github.com/wb-go/wbf/zlog"
)

type Service struct {
	NoticeService *noticeservice.NoticeService
	rp            *repository.Repository
}

func New(rp *repository.Repository, rb *rabbitmq.Client) *Service {
	lgNoticeService := zlog.Logger.With().Str("component", "service-noticeService").Logger()

	return &Service{
		NoticeService: noticeservice.New(lgNoticeService, rb),
		rp:            rp,
	}
}
