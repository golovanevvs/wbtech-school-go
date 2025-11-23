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

// EventRepository implements the event repository interface for PostgreSQL
type EventRepository struct {
	db *pkgPostgres.Postgres
}

// NewEventRepository creates a new instance of EventRepository
func NewEventRepository(db *pkgPostgres.Postgres) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event
func (rp *EventRepository) Create(event *model.Event) (*model.Event, error) {
	query := `

		INSERT INTO
			events (title, date, description, total_places, available_places, booking_deadline, owner_id, created_at, updated_at) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING
			id, title, date, description, total_places, available_places, booking_deadline, owner_id, created_at, updated_at
		
		`

	var createdEvent model.Event
	err := rp.db.DB.Master.QueryRowContext(
		context.Background(),
		query,
		event.Title,
		event.Date,
		event.Description,
		event.TotalPlaces,
		event.AvailablePlaces,
		event.BookingDeadline,
		event.OwnerID,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(
		&createdEvent.ID,
		&createdEvent.Title,
		&createdEvent.Date,
		&createdEvent.Description,
		&createdEvent.TotalPlaces,
		&createdEvent.AvailablePlaces,
		&createdEvent.BookingDeadline,
		&createdEvent.OwnerID,
		&createdEvent.CreatedAt,
		&createdEvent.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &createdEvent, nil
}

// GetByID returns an event by ID
func (rp *EventRepository) GetByID(id int) (*model.Event, error) {
	query := `

		SELECT
			id, title, date, description, total_places, available_places, booking_deadline, owner_id, created_at, updated_at
		FROM
			events
		WHERE
			id = $1
		
		`
	var event model.Event

	row := rp.db.DB.QueryRowContext(context.Background(), query, id)
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Date,
		&event.Description,
		&event.TotalPlaces,
		&event.AvailablePlaces,
		&event.BookingDeadline,
		&event.OwnerID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// GetAll returns all events
func (rp *EventRepository) GetAll() ([]*model.Event, error) {
	query := `

		SELECT
			id, title, date, description, total_places, available_places, booking_deadline, owner_id, created_at, updated_at
		FROM
			events
		ORDER BY
			date
		
		`

	rows, err := rp.db.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Description,
			&event.TotalPlaces,
			&event.AvailablePlaces,
			&event.BookingDeadline,
			&event.OwnerID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}

// Update updates an event
func (rp *EventRepository) Update(event *model.Event) error {
	// Если изменяется количество мест, нужно пересчитать available_places
	if event.TotalPlaces != 0 {
		// Получаем текущее состояние события
		currentEvent, err := rp.GetByID(event.ID)
		if err != nil {
			return fmt.Errorf("failed to get current event for available places calculation: %w", err)
		}

		// Вычисляем количество занятых мест
		bookedPlaces := currentEvent.TotalPlaces - currentEvent.AvailablePlaces

		// Пересчитываем доступные места
		event.AvailablePlaces = event.TotalPlaces - bookedPlaces

		// Если доступных мест получилось отрицательное число, устанавливаем в 0
		if event.AvailablePlaces < 0 {
			event.AvailablePlaces = 0
		}
	}

	query := `

		UPDATE
			events 
		SET
			title = $1, date = $2, description = $3, total_places = $4, available_places = $5, booking_deadline = $6, updated_at = $7 
		WHERE
			id = $8
		
		`

	result, err := rp.db.DB.ExecContext(
		context.Background(),
		query,
		event.Title,
		event.Date,
		event.Description,
		event.TotalPlaces,
		event.AvailablePlaces,
		event.BookingDeadline,
		event.UpdatedAt,
		event.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event with id %d not found", event.ID)
	}

	return nil
}

// Delete deletes an event by ID
func (rp *EventRepository) Delete(id int) error {
	query := `

		DELETE FROM
			events
		WHERE
			id = $1
		
		`
	result, err := rp.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event with id %d not found", id)
	}

	return nil
}

// UpdateAvailablePlaces updates the available places for an event
func (rp *EventRepository) UpdateAvailablePlaces(eventID int, newAvailablePlaces int) error {
	query := `

		UPDATE
			events
		SET
			available_places = $1, updated_at = $2
		WHERE
			id = $3
		
		`
	currentTime := time.Now()

	result, err := rp.db.DB.ExecContext(
		context.Background(),
		query,
		newAvailablePlaces,
		currentTime,
		eventID,
	)
	if err != nil {
		return fmt.Errorf("failed to update available places for event %d: %w", eventID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event with id %d not found", eventID)
	}

	return nil
}

// GetByOwnerID returns all events owned by a user
func (rp *EventRepository) GetByOwnerID(ownerID int) ([]*model.Event, error) {
	query := `

		SELECT
			id, title, date, description, total_places, available_places, booking_deadline, owner_id, created_at, updated_at
		FROM
			events
		WHERE
			owner_id = $1
		ORDER BY
			date
		
		`

	rows, err := rp.db.DB.QueryContext(context.Background(), query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by owner ID: %w", err)
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Date,
			&event.Description,
			&event.TotalPlaces,
			&event.AvailablePlaces,
			&event.BookingDeadline,
			&event.OwnerID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
}
