package service

import "github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/internal/repository"

type Service struct {
	rp *repository.Repository
}

func New(rp *repository.Repository) *Service {
	return &Service{
		rp: rp,
	}
}
