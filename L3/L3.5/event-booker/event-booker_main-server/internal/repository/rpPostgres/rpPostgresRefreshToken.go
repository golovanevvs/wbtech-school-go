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

// RefreshTokenRepository implements the refresh token repository interface for PostgreSQL
type RefreshTokenRepository struct {
	db *pkgPostgres.Postgres
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
func NewRefreshTokenRepository(db *pkgPostgres.Postgres) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(token *model.RefreshToken) (*model.RefreshToken, error) {
	query := `

		INSERT INTO
			refresh_tokens (user_id, token, expires_at, created_at) 
		VALUES
			($1, $2, $3, $4) 
		RETURNING
			id, user_id, token, expires_at, created_at
		
		`

	var createdToken model.RefreshToken
	err := r.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(
		&createdToken.ID,
		&createdToken.UserID,
		&createdToken.Token,
		&createdToken.ExpiresAt,
		&createdToken.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &createdToken, nil
}

// GetByToken returns a refresh token by its value
func (r *RefreshTokenRepository) GetByToken(token string) (*model.RefreshToken, error) {
	query := `

		SELECT
			id, user_id, token, expires_at, created_at
		FROM
			refresh_tokens
		WHERE
			token = $1
		
		`
	var refreshToken model.RefreshToken

	row := r.db.DB.QueryRowContext(context.Background(), query, token)
	err := row.Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &refreshToken, nil
}

// GetByUserID returns all refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(userID int) ([]*model.RefreshToken, error) {
	query := `

		SELECT
			id, user_id, token, expires_at, created_at
		FROM
			refresh_tokens
		WHERE
			user_id = $1
		ORDER BY
			created_at DESC
		
		`

	rows, err := r.db.DB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens for user %d: %w", userID, err)
	}
	defer rows.Close()

	var tokens []*model.RefreshToken
	for rows.Next() {
		var token model.RefreshToken
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refresh token: %w", err)
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

// DeleteByID deletes a refresh token by ID
func (r *RefreshTokenRepository) DeleteByID(id int) error {
	query := `
	
		DELETE FROM
			refresh_tokens
		WHERE
			id = $1
		
		`
	result, err := r.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token by ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token with id %d not found", id)
	}

	return nil
}

// DeleteByToken deletes a refresh token by its value
func (r *RefreshTokenRepository) DeleteByToken(token string) error {
	query := `
	
		DELETE FROM
			refresh_tokens
		WHERE
			token = $1
		
		`
	result, err := r.db.DB.ExecContext(context.Background(), query, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token by token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token with value %s not found", token)
	}

	return nil
}

// DeleteByUserID deletes all refresh tokens for a user
func (r *RefreshTokenRepository) DeleteByUserID(userID int) error {
	query := `
	
		DELETE FROM
			refresh_tokens
		WHERE
			user_id = $1
		
		`
	result, err := r.db.DB.ExecContext(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens for user %d: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no refresh tokens found for user %d", userID)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired() error {
	query := `
	
		DELETE FROM
			refresh_tokens
		WHERE
			expires_at < $1
		
		`
	currentTime := time.Now()

	result, err := r.db.DB.ExecContext(context.Background(), query, currentTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	return nil
}
