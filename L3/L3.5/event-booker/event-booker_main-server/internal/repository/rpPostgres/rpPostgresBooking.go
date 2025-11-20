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

// BookingRepository implements the booking repository interface for PostgreSQL
type BookingRepository struct {
	db *pkgPostgres.Postgres
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *pkgPostgres.Postgres) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create creates a new booking
func (r *BookingRepository) Create(booking *model.Booking) (*model.Booking, error) {
	query := `

		INSERT INTO
			bookings (user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7) 
		RETURNING
			id, user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at
		
		`

	var createdBooking model.Booking
	err := r.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		booking.UserID,
		booking.EventID,
		booking.Status,
		booking.CreatedAt,
		booking.ExpiresAt,
		booking.ConfirmedAt,
		booking.CancelledAt,
	).Scan(
		&createdBooking.ID,
		&createdBooking.UserID,
		&createdBooking.EventID,
		&createdBooking.Status,
		&createdBooking.CreatedAt,
		&createdBooking.ExpiresAt,
		&createdBooking.ConfirmedAt,
		&createdBooking.CancelledAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return &createdBooking, nil
}

// GetByID returns a booking by ID
func (r *BookingRepository) GetByID(id int) (*model.Booking, error) {
	query := `

		SELECT
			id, user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at
		FROM
			bookings
		WHERE
			id = $1
		
		`
	var booking model.Booking

	row := r.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.EventID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ExpiresAt,
		&booking.ConfirmedAt,
		&booking.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("booking with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// GetByUserID returns all bookings for a user
func (r *BookingRepository) GetByUserID(userID int) ([]*model.Booking, error) {
	query := `

		SELECT
			id, user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at
		FROM
			bookings
		WHERE
			user_id = $1
		ORDER BY
			created_at DESC
		
		`

	rows, err := r.db.DB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings for user %d: %w", userID, err)
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.EventID,
			&booking.Status,
			&booking.CreatedAt,
			&booking.ExpiresAt,
			&booking.ConfirmedAt,
			&booking.CancelledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

// GetByEventID returns all bookings for an event
func (r *BookingRepository) GetByEventID(eventID int) ([]*model.Booking, error) {
	query := `

		SELECT
			id, user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at
		FROM
			bookings
		WHERE
			event_id = $1
		ORDER BY
			created_at DESC
		
		`

	rows, err := r.db.DB.QueryContext(context.Background(), query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings for event %d: %w", eventID, err)
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.EventID,
			&booking.Status,
			&booking.CreatedAt,
			&booking.ExpiresAt,
			&booking.ConfirmedAt,
			&booking.CancelledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

// Update updates a booking
func (r *BookingRepository) Update(booking *model.Booking) error {
	query := `

		UPDATE
			bookings 
		SET
			user_id = $1, event_id = $2, status = $3, expires_at = $4, confirmed_at = $5, cancelled_at = $6 
		WHERE
			id = $7
		
		`

	result, err := r.db.DB.ExecContext(
		context.Background(),
		query,
		booking.UserID,
		booking.EventID,
		booking.Status,
		booking.ExpiresAt,
		booking.ConfirmedAt,
		booking.CancelledAt,
		booking.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking with id %d not found", booking.ID)
	}

	return nil
}

// Delete deletes a booking by ID
func (r *BookingRepository) Delete(id int) error {
	query := `
	
		DELETE FROM
			bookings
		WHERE
			id = $1
		
		`
	result, err := r.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking with id %d not found", id)
	}

	return nil
}

// GetExpiredBookings returns all bookings that have expired
func (r *BookingRepository) GetExpiredBookings() ([]*model.Booking, error) {
	query := `

		SELECT
			id, user_id, event_id, status, created_at, expires_at, confirmed_at, cancelled_at
		FROM
			bookings
		WHERE
			expires_at < $1 AND status = $2
		
		`

	currentTime := time.Now()
	rows, err := r.db.DB.QueryContext(context.Background(), query, currentTime, model.BookingPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.EventID,
			&booking.Status,
			&booking.CreatedAt,
			&booking.ExpiresAt,
			&booking.ConfirmedAt,
			&booking.CancelledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expired booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

// UpdateStatus updates the status of a booking
func (r *BookingRepository) UpdateStatus(id int, status model.BookingStatus) error {
	query := `

		UPDATE
			bookings
		SET
			status = $1, updated_at = $2
		WHERE
			id = $3
		
		`
	currentTime := time.Now()

	result, err := r.db.DB.ExecContext(
		context.Background(),
		query,
		status,
		currentTime,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update booking status for booking %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking with id %d not found", id)
	}

	return nil
}
