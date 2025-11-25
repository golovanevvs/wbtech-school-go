package service

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/repository"

// Service structure that combines all services
type Service struct {
	salesRp ISalesRepository
}

// New creates a new Service structure
func New(rp *repository.Repository) *Service {
	return &Service{
		salesRp: rp,
	}
}
