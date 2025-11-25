package rpPostgres

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
)

// CreateSalesRecord creates a new sales record in the database
func (rp *RpPostgres) CreateSalesRecord(ctx context.Context, data model.Data) (int, error) {
	amountDB := data.Amount * 100

	query := `
	
		INSERT INTO sales_records
			(type, category, date, amount)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id
	
	`

	var id int
	row := rp.db.DB.QueryRowContext(ctx, query, data.Type, data.Category, data.Date, amountDB)
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create sales record in Postgres: %w", err)
	}

	return id, nil
}
