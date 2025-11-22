package rpPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
func (rp *UserRepository) Create(user *model.User) (*model.User, error) {
	query := `

		INSERT INTO users
			(email, name, password_hash, telegram_notifications, email_notifications, created_at, updated_at) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7) 
		RETURNING
			id
		
		`

	var createdUser model.User
	createdUser = *user
	err := rp.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.TelegramNotifications,
		user.EmailNotifications,
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
			id, email, name, password_hash, telegram_username, telegram_chat_id, telegram_notifications, email_notifications, created_at, updated_at 
		FROM
			users
		WHERE
			id = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.TelegramNotifications,
		&user.EmailNotifications,
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

// GetByEmail returns a user by email
func (rp *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `

		SELECT
			id, email, name, password_hash, telegram_username, telegram_chat_id, telegram_notifications, email_notifications, created_at, updated_at 
		FROM
			users
		WHERE
			email = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.TelegramNotifications,
		&user.EmailNotifications,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
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
            telegram_username = $2, 
            telegram_notifications = $3,
            email_notifications = $4,
            updated_at = $5  
		WHERE
			id = $6
		
		`

	result, err := rp.db.DB.ExecContext(
		context.Background(),
		query,
		user.Name,
		user.TelegramUsername,
		user.TelegramNotifications,
		user.EmailNotifications,
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

// SaveTelegramChatID saves the Telegram chat ID for a user
func (rp *UserRepository) SaveTelegramChatID(ctx context.Context, userID int, chatID int64) error {
	query := `

		UPDATE
			users 
		SET
			telegram_chat_id = $1, updated_at = $2 
		WHERE
			id = $3
		
		`

	result, err := rp.db.DB.ExecContext(ctx, query, &chatID, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update telegram chat ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

// GetByTelegramChatID returns a user by Telegram chat ID
func (rp *UserRepository) GetByTelegramChatID(ctx context.Context, chatID int64) (*model.User, error) {
	query := `
		
		SELECT
			id, email, name, password_hash, telegram_username, telegram_chat_id, telegram_notifications, email_notifications, created_at, updated_at 
		FROM
			users
		WHERE
			telegram_chat_id = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, chatID)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.TelegramNotifications,
		&user.EmailNotifications,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with telegram chat id %d not found", chatID)
		}
		return nil, fmt.Errorf("failed to get user by telegram chat id: %w", err)
	}

	return &user, nil
}

// UpdateTelegramUsername updates the Telegram username for a user
func (rp *UserRepository) UpdateTelegramUsername(ctx context.Context, userID int, username *string) error {
	query := `

		UPDATE
			users 
		SET
			telegram_username = $1, updated_at = $2 
		WHERE
			id = $3
		
		`

	result, err := rp.db.DB.ExecContext(ctx, query, username, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update telegram username: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

// GetByTelegramUsername returns a user by Telegram username
func (rp *UserRepository) GetByTelegramUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		
		SELECT
			id, email, name, password_hash, telegram_username, telegram_chat_id, telegram_notifications, email_notifications, created_at, updated_at 
		FROM
			users
		WHERE
			telegram_username = $1
		
		`
	var user model.User

	row := rp.db.DB.QueryRowContext(context.Background(), query, username)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.TelegramUsername,
		&user.TelegramChatID,
		&user.TelegramNotifications,
		&user.EmailNotifications,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with telegram username %s not found", username)
		}
		return nil, fmt.Errorf("failed to get user by telegram username: %w", err)
	}

	return &user, nil
}
