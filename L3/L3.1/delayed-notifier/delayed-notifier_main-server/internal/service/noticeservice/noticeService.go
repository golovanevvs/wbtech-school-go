package noticeservice

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type NoticeService struct {
	lg zlog.Zerolog
	rb *rabbitmq.Client
}

func New(lg zlog.Zerolog, rb *rabbitmq.Client) *NoticeService {
	return &NoticeService{
		lg: lg,
		rb: rb,
	}
}
