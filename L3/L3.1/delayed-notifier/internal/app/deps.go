package app

import (
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/transport"
)

type dependencies struct {
	rm *resourceManager

	tr *transport.Transport
	rp *repository.Repository
	sv *service.Service
}

func newDependenies(cfg *appConfig) *dependencies {
	tr := transport.New
	return &dependencies{}
}
