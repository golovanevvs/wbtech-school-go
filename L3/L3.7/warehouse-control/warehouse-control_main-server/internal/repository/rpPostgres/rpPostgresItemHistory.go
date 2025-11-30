package rpPostgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"
)

// ItemHistoryRepository implements the item history repository interface for PostgreSQL
type ItemHistoryRepository struct {
	db *pkgPostgres.Postgres
}

// NewItemHistoryRepository creates a new instance of ItemHistoryRepository
func NewItemHistoryRepository(db *pkgPostgres.Postgres) *ItemHistoryRepository {
	return &ItemHistoryRepository{db: db}
}

// GetByItemID returns all history records for a specific item
func (rp *ItemHistoryRepository) GetByItemID(itemID int) ([]model.ItemAction, error) {
	query := `
		SELECT
			id, item_id, action_type, user_id, user_name, changes, created_at
		FROM
			item_actions
		WHERE
			item_id = $1
		ORDER BY
			created_at DESC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query item history: %w", err)
	}
	defer rows.Close()

	var history []model.ItemAction
	for rows.Next() {
		var action model.ItemAction
		var changesJSON sql.NullString
		err := rows.Scan(
			&action.ID,
			&action.ItemID,
			&action.ActionType,
			&action.UserID,
			&action.UserName,
			&changesJSON,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item action: %w", err)
		}

		if changesJSON.Valid {
			action.Changes = changesJSON.String
		} else {
			action.Changes = "null"
		}

		history = append(history, action)
	}

	return history, nil
}

// GetAll returns all item history records
func (rp *ItemHistoryRepository) GetAll() ([]model.ItemAction, error) {
	query := `
		SELECT
			id, item_id, action_type, user_id, user_name, changes, created_at
		FROM
			item_actions
		ORDER BY
			created_at DESC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all item history: %w", err)
	}
	defer rows.Close()

	var history []model.ItemAction
	for rows.Next() {
		var action model.ItemAction
		var changesJSON sql.NullString
		err := rows.Scan(
			&action.ID,
			&action.ItemID,
			&action.ActionType,
			&action.UserID,
			&action.UserName,
			&changesJSON,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item action: %w", err)
		}

		if changesJSON.Valid {
			action.Changes = changesJSON.String
		} else {
			action.Changes = "null"
		}

		history = append(history, action)
	}

	return history, nil
}

// CreateAction creates a new item action record
func (rp *ItemHistoryRepository) CreateAction(itemID int, actionType string, userID int, userName string, changes map[string]interface{}) error {
	query := `
		INSERT INTO item_actions (item_id, action_type, user_id, user_name, changes)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := rp.db.DB.ExecContext(context.Background(), query, itemID, actionType, userID, userName, changes)
	if err != nil {
		return fmt.Errorf("failed to create item action: %w", err)
	}

	return nil
}

// ExportToCSV returns history data formatted for CSV export
func (rp *ItemHistoryRepository) ExportToCSV(itemID int) ([]map[string]interface{}, error) {
	query := `
		SELECT
			ia.id,
			ia.action_type,
			ia.user_name,
			ia.created_at,
			ia.changes,
			i.name as item_name,
			i.price as item_price,
			i.quantity as item_quantity
		FROM
			item_actions ia
		JOIN
			items i ON ia.item_id = i.id
		WHERE
			ia.item_id = $1
		ORDER BY
			ia.created_at DESC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query item history for export: %w", err)
	}
	defer rows.Close()

	var csvData []map[string]interface{}
	for rows.Next() {
		var actionType, userName, itemName string
		var actionID int
		var itemPrice string // Изменено с int на string для DECIMAL
		var itemQuantity int
		var createdAt string
		var changesJSON sql.NullString

		err := rows.Scan(
			&actionID,
			&actionType,
			&userName,
			&createdAt,
			&changesJSON,
			&itemName,
			&itemPrice,
			&itemQuantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row for export: %w", err)
		}

		row := map[string]interface{}{
			"ID":           actionID,
			"Товар":        itemName,
			"Действие":     actionType,
			"Пользователь": userName,
			"Дата":         createdAt,
			"Изменения":    changesJSON.String,
		}

		csvData = append(csvData, row)
	}

	return csvData, nil
}
