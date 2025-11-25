package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
)

// ISalesRepository defines the interface for SalesRecord database operations
type ISalesRepository interface {
	CreateSalesRecord(ctx context.Context, data model.Data) (int, error)
}

// CreateSalesRecord creates a new sales record
func (sv *Service) CreateSalesRecord(ctx context.Context, data model.Data) (int, error) {
	return sv.salesRp.CreateSalesRecord(ctx, data)
}
