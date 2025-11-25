package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
)

// ISalesRepository defines the interface for SalesRecord database operations
type ISalesRepository interface {
	CreateSalesRecord(ctx context.Context, data model.Data) (int, error)
	GetSalesRecords(ctx context.Context, sortOptions model.SortOptions) ([]model.Data, error)
	UpdateSalesRecord(ctx context.Context, id int, data model.Data) error
	DeleteSalesRecord(ctx context.Context, id int) error
}

// CreateSalesRecord creates a new sales record
func (sv *Service) CreateSalesRecord(ctx context.Context, data model.Data) (int, error) {
	return sv.salesRp.CreateSalesRecord(ctx, data)
}

// GetSalesRecords retrieves sales records with sorting
func (sv *Service) GetSalesRecords(ctx context.Context, sortOptions model.SortOptions) ([]model.Data, error) {
	return sv.salesRp.GetSalesRecords(ctx, sortOptions)
}

// UpdateSalesRecord updates an existing sales record
func (sv *Service) UpdateSalesRecord(ctx context.Context, id int, data model.Data) error {
	return sv.salesRp.UpdateSalesRecord(ctx, id, data)
}

// DeleteSalesRecord deletes a sales record
func (sv *Service) DeleteSalesRecord(ctx context.Context, id int) error {
	return sv.salesRp.DeleteSalesRecord(ctx, id)
}
