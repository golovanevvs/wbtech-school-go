package rpPostgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	orderBy := "id"
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
		SELECT
			id, type, category, date, amount
		FROM
			sales_records
		ORDER BY
			%s %s
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
		UPDATE
			sales_records
		SET
			type = $1, category = $2, date = $3, amount = $4
		WHERE
			id = $5
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
	query := `
		SELECT 
			SUM(amount) as total_sum,
			AVG(amount) as avg_amount,
			COUNT(*) as record_count,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) as median_amount,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) as percentile_90_amount
		FROM
			sales_records
		WHERE
			date >= $1 AND date <= $2
	`

	rows, err := rp.db.DB.QueryContext(ctx, query, from, to)
	if err != nil {
		return model.Analytics{}, fmt.Errorf("failed to get analytics from Postgres: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Analytics{
			Sum:          0,
			Avg:          0,
			Count:        0,
			Median:       0,
			Percentile90: 0,
		}, nil
	}

	var totalSum, avgAmount, medianAmount, percentile90Amount sql.NullFloat64
	var recordCount int

	err = rows.Scan(&totalSum, &avgAmount, &recordCount, &medianAmount, &percentile90Amount)
	if err != nil {
		return model.Analytics{}, fmt.Errorf("failed to scan analytics data: %w", err)
	}

	sum := 0.0
	if totalSum.Valid {
		sum = totalSum.Float64 / 100
	}

	avg := 0.0
	if avgAmount.Valid {
		avg = avgAmount.Float64 / 100
	}

	median := 0.0
	if medianAmount.Valid {
		median = medianAmount.Float64 / 100
	}

	percentile90 := 0.0
	if percentile90Amount.Valid {
		percentile90 = percentile90Amount.Float64 / 100
	}

	return model.Analytics{
		Sum:          sum,
		Avg:          avg,
		Count:        recordCount,
		Median:       median,
		Percentile90: percentile90,
	}, nil
}

// ExportCSV exports sales records to CSV format
func (rp *RpPostgres) ExportCSV(ctx context.Context, from, to string) ([]byte, error) {
	query := `
		SELECT
			id, type, category, date, amount
		FROM
			sales_records
	`

	var args []interface{}
	var whereClause string

	if from != "" {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "date >= $" + fmt.Sprint(len(args)+1)
		args = append(args, from)
	}

	if to != "" {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "date <= $" + fmt.Sprint(len(args)+1)
		args = append(args, to)
	}

	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	query += " ORDER BY date DESC"

	rows, err := rp.db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales records for CSV export: %w", err)
	}
	defer rows.Close()

	var csvContent strings.Builder

	csvContent.WriteString("ID,Type,Category,Date,Amount\n")

	for rows.Next() {
		var record model.Data
		var amountDB int

		err := rows.Scan(&record.ID, &record.Type, &record.Category, &record.Date, &amountDB)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales record for CSV: %w", err)
		}

		record.Amount = float64(amountDB) / 100

		escapeCSVField := func(field string) string {
			if strings.Contains(field, ",") || strings.Contains(field, "\"") || strings.Contains(field, "\n") {
				return "\"" + strings.ReplaceAll(field, "\"", "\"\"") + "\""
			}
			return field
		}

		csvContent.WriteString(fmt.Sprintf("%d,%s,%s,%s,%.2f\n",
			record.ID,
			escapeCSVField(record.Type),
			escapeCSVField(record.Category),
			escapeCSVField(record.Date),
			record.Amount,
		))
	}

	return []byte(csvContent.String()), nil
}

// DeleteSalesRecord deletes a sales record
func (rp *RpPostgres) DeleteSalesRecord(ctx context.Context, id int) error {
	query := `
		DELETE FROM
			sales_records
		WHERE
			id = $1
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
