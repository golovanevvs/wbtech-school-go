package rpPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"
)

// UserRepository implements the user repository interface for PostgreSQL
type UserRepository struct {
	db *pkgPostgres.Postgres
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *pkgPostgres.Postgres) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (rp *UserRepository) Create(user *model.User) (*model.User, error) {
	query := `

		INSERT INTO users
			(username, name, password_hash, user_role, created_at, updated_at) 
		VALUES
			($1, $2, $3, $4, $5, $6) 
		RETURNING
			id
		
		`

	var createdUser model.User
	createdUser = *user
	err := rp.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		user.UserName,
		user.Name,
		user.PasswordHash,
		user.UserRole,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&createdUser.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &createdUser, nil
}

// GetByID returns a user by ID
func (rp *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
	
		SELECT
			id, username, name, password_hash, user_role, created_at, updated_at 
		FROM
			users
		WHERE
			id = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Name,
		&user.PasswordHash,
		&user.UserRole,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByUsername returns a user by username
func (rp *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `

		SELECT
			id, username, name, password_hash, user_role, created_at, updated_at 
		FROM
			users
		WHERE
			username = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, username)
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Name,
		&user.PasswordHash,
		&user.UserRole,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (rp *UserRepository) Update(user *model.User) error {
	query := `

		UPDATE
			users 
		SET
			name = $1, 
            user_role = $2,
            updated_at = $3  
		WHERE
			id = $4
		
		`

	result, err := rp.db.DB.ExecContext(
		context.Background(),
		query,
		user.Name,
		user.UserRole,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return nil
}

// Delete deletes a user by ID
func (rp *UserRepository) Delete(id int) error {
	query := `
	
		DELETE FROM
			users
		WHERE
			id = $1
		
		`
	result, err := rp.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
