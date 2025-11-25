package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/repository"
	"github.com/wb-go/wbf/zlog"
)

// Service structure that combines all services
type Service struct {
}

// New creates a new Service structure
func New(
	cfg *Config,
	rp *repository.Repository,
) *Service {
	lg := zlog.Logger.With().Str("layer", "service").Logger()

	return &Service{}
}
