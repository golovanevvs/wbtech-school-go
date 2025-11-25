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
