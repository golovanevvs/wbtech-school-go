package rpPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"
)

// ItemRepository implements the item repository interface for PostgreSQL
type ItemRepository struct {
	db *pkgPostgres.Postgres
}

// NewItemRepository creates a new instance of ItemRepository
func NewItemRepository(db *pkgPostgres.Postgres) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create creates a new item
func (rp *ItemRepository) Create(item *model.Item, userID int, userName string) (*model.Item, error) {
	// Устанавливаем контекст пользователя для триггера
	_, err := rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_id = "+strconv.Itoa(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to set user_id: %w", err)
	}

	_, err = rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_name = '"+userName+"'")
	if err != nil {
		return nil, fmt.Errorf("failed to set user_name: %w", err)
	}

	query := `
		INSERT INTO items
			(name, price, quantity, created_at, updated_at) 
		VALUES
			($1, $2, $3, $4, $5) 
		RETURNING
			id
	`

	var createdItem model.Item
	createdItem = *item
	err = rp.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		item.Name,
		item.Price,
		item.Quantity,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&createdItem.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &createdItem, nil
}

// GetAll returns all items
func (rp *ItemRepository) GetAll() ([]model.Item, error) {
	query := `
		SELECT
			id, name, price, quantity, created_at, updated_at 
		FROM
			items 
		ORDER BY
			id DESC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetByID returns an item by ID
func (rp *ItemRepository) GetByID(id int) (*model.Item, error) {
	query := `
		SELECT
			id, name, price, quantity, created_at, updated_at 
		FROM
			items
		WHERE
			id = $1
	`

	var item model.Item
	row := rp.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

// Update updates an item
func (rp *ItemRepository) Update(item *model.Item, userID int, userName string) error {
	// Устанавливаем контекст пользователя для триггера
	_, err := rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_id = "+strconv.Itoa(userID))
	if err != nil {
		return fmt.Errorf("failed to set user_id: %w", err)
	}

	_, err = rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_name = '"+userName+"'")
	if err != nil {
		return fmt.Errorf("failed to set user_name: %w", err)
	}

	query := `
		UPDATE
			items 
		SET
			name = $1, 
			price = $2,
			quantity = $3,
			updated_at = $4  
		WHERE
			id = $5
	`

	result, err := rp.db.DB.ExecContext(
		context.Background(),
		query,
		item.Name,
		item.Price,
		item.Quantity,
		item.UpdatedAt,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with id %d not found", item.ID)
	}

	return nil
}

// Delete deletes an item by ID
func (rp *ItemRepository) Delete(id int, userID int, userName string) error {
	// Устанавливаем контекст пользователя для триггера
	_, err := rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_id = "+strconv.Itoa(userID))
	if err != nil {
		return fmt.Errorf("failed to set user_id: %w", err)
	}

	_, err = rp.db.DB.Master.ExecContext(context.Background(), "SET app.current_user_name = '"+userName+"'")
	if err != nil {
		return fmt.Errorf("failed to set user_name: %w", err)
	}

	query := `
		DELETE FROM
			items
		WHERE
			id = $1
	`
	result, err := rp.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with id %d not found", id)
	}

	return nil
}
