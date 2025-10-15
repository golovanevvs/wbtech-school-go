package publishNoticeService

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
)

type PublishNoticeService struct {
	rb *pkgRabbitmq.Client
}

func New(rb *pkgRabbitmq.Client) *PublishNoticeService {
	return &PublishNoticeService{
		rb: rb,
	}
}

func (sv *PublishNoticeService) Publish(ctx context.Context, notice model.Notice) (err error) {
	return nil
}
