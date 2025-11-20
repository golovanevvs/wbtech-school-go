package rpPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgPostgres"
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
func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	query := `
		
		INSERT INTO users
			(email, name, created_at, updated_at) 
		VALUES
			($1, $2, $3, $4) 
		RETURNING
			id, email, name, created_at, updated_at
		
		`

	var createdUser model.User
	err := r.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Name,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &createdUser, nil
}

// GetByID returns a user by ID
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
	
		SELECT
			id, email, name, created_at, updated_at
		FROM 
			users
		WHERE
			id = $1
		
		`

	var user model.User

	row := r.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail returns a user by email
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `

		SELECT
			id, email, name, created_at, updated_at
		FROM
			users
		WHERE
			email = $1
		
		`

	var user model.User

	row := r.db.DB.QueryRowContext(context.Background(), query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	query := `

		UPDATE
			users 
		SET
			email = $1, name = $2, updated_at = $3 
		WHERE
			id = $4
		
		`

	result, err := r.db.DB.ExecContext(
		context.Background(),
		query,
		user.Email,
		user.Name,
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
func (r *UserRepository) Delete(id int) error {
	query := `

		DELETE FROM
			users
		WHERE
			id = $1
		
		`

	result, err := r.db.DB.ExecContext(context.Background(), query, id)
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
