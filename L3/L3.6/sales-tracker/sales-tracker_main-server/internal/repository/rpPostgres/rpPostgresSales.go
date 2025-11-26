package rpPostgres

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
)

// CreateSalesRecord creates a new sales record in the database
func (rp *RpPostgres) CreateSalesRecord(ctx context.Context, data model.Data) (int, error) {
	amountDB := int(data.Amount * 100)

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

// GetSalesRecords retrieves sales records with sorting
func (rp *RpPostgres) GetSalesRecords(ctx context.Context, sortOptions model.SortOptions) ([]model.Data, error) {
	// Build the ORDER BY clause based on sort options
	orderBy := "id" // default sorting by ID
	if sortOptions.Field != "" {
		switch sortOptions.Field {
		case "id":
			orderBy = "id"
		case "type":
			orderBy = "type"
		case "category":
			orderBy = "category"
		case "date":
			orderBy = "date"
		case "amount":
			orderBy = "amount"
		default:
			orderBy = "id"
		}
	}

	direction := "ASC"
	if sortOptions.Direction != "" {
		if sortOptions.Direction == "desc" || sortOptions.Direction == "DESC" {
			direction = "DESC"
		}
	}

	query := fmt.Sprintf(`
		SELECT id, type, category, date, amount
		FROM sales_records
		ORDER BY %s %s
	`, orderBy, direction)

	rows, err := rp.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales records from Postgres: %w", err)
	}
	defer rows.Close()

	var records []model.Data
	for rows.Next() {
		var record model.Data
		var amountDB int
		err := rows.Scan(&record.ID, &record.Type, &record.Category, &record.Date, &amountDB)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales record: %w", err)
		}
		record.Amount = float64(amountDB) / 100
		records = append(records, record)
	}

	return records, nil
}

// UpdateSalesRecord updates an existing sales record
func (rp *RpPostgres) UpdateSalesRecord(ctx context.Context, id int, data model.Data) error {
	amountDB := int(data.Amount * 100)

	query := `
		UPDATE sales_records
		SET type = $1, category = $2, date = $3, amount = $4
		WHERE id = $5
	`

	result, err := rp.db.DB.ExecContext(ctx, query, data.Type, data.Category, data.Date, amountDB, id)
	if err != nil {
		return fmt.Errorf("failed to update sales record in Postgres: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sales record with id %d not found", id)
	}

	return nil
}

// GetAnalytics retrieves analytics data for a given period
func (rp *RpPostgres) GetAnalytics(ctx context.Context, from, to string) (model.Analytics, error) {
	// First, get all records within the date range
	query := `
		SELECT amount
		FROM sales_records
		WHERE date >= $1 AND date <= $2
		ORDER BY amount
	`

	rows, err := rp.db.DB.QueryContext(ctx, query, from, to)
	if err != nil {
		return model.Analytics{}, fmt.Errorf("failed to get sales records for analytics: %w", err)
	}
	defer rows.Close()

	var amounts []float64
	for rows.Next() {
		var amountDB int
		err := rows.Scan(&amountDB)
		if err != nil {
			return model.Analytics{}, fmt.Errorf("failed to scan amount: %w", err)
		}
		amounts = append(amounts, float64(amountDB)/100)
	}

	// If no records found, return empty analytics
	if len(amounts) == 0 {
		return model.Analytics{
			Sum:          0,
			Avg:          0,
			Count:        0,
			Median:       0,
			Percentile90: 0,
		}, nil
	}

	// Calculate basic metrics
	var sum float64
	for _, amount := range amounts {
		sum += amount
	}
	count := len(amounts)
	avg := sum / float64(count)

	// Calculate median
	var median float64
	if count%2 == 0 {
		// Even number of elements
		median = (amounts[count/2-1] + amounts[count/2]) / 2
	} else {
		// Odd number of elements
		median = amounts[count/2]
	}

	// Calculate 90th percentile
	percentile90Index := int(float64(count) * 0.9)
	if percentile90Index >= count {
		percentile90Index = count - 1
	}
	if percentile90Index < 0 {
		percentile90Index = 0
	}
	percentile90 := amounts[percentile90Index]

	return model.Analytics{
		Sum:          sum,
		Avg:          avg,
		Count:        count,
		Median:       median,
		Percentile90: percentile90,
	}, nil
}

// DeleteSalesRecord deletes a sales record
func (rp *RpPostgres) DeleteSalesRecord(ctx context.Context, id int) error {
	query := `
		DELETE FROM sales_records
		WHERE id = $1
	`

	result, err := rp.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sales record from Postgres: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sales record with id %d not found", id)
	}

	return nil
}
