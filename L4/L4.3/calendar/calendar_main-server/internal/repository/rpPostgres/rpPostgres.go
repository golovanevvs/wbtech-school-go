package rpPostgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgPostgres"
)

// RpPostgres implements the repository interface for PostgreSQL
type RpPostgres struct {
	db *pkgPostgres.Postgres
}

// NewPostgresRepository creates a new instance of PostgresRepository
func New(db *pkgPostgres.Postgres) *RpPostgres {
	return &RpPostgres{db: db}
}

// Close closes the database connection
func (rp *RpPostgres) Close() error {
	return rp.db.Close()
}

// GetMonthEvents returns events for a specific month
func (rp *RpPostgres) GetMonthEvents(year, month int) ([]model.Event, error) {
	query := `
		SELECT id, title, description, start, end, all_day, reminder, reminder_time, created_at, updated_at
		FROM events
		WHERE EXTRACT(YEAR FROM start) = $1 AND EXTRACT(MONTH FROM start) = $2
		ORDER BY start ASC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query month events: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		var reminderTime sql.NullTime
		var endTime sql.NullTime

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Start,
			&endTime,
			&event.AllDay,
			&event.Reminder,
			&reminderTime,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		if endTime.Valid {
			event.End = &endTime.Time
		}
		if reminderTime.Valid {
			event.ReminderTime = &reminderTime.Time
		}

		events = append(events, event)
	}

	return events, nil
}

// CreateEvent creates a new event
func (rp *RpPostgres) CreateEvent(eventData *model.CreateEventRequest) (*model.Event, error) {
	query := `
		INSERT INTO events (title, description, start, end, all_day, reminder, reminder_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, start, end, all_day, reminder, reminder_time, created_at, updated_at
	`

	var event model.Event
	var endTime sql.NullTime
	var reminderTime sql.NullTime

	err := rp.db.DB.QueryRowContext(context.Background(),
		query,
		eventData.Title,
		eventData.Description,
		eventData.Start,
		eventData.End,
		eventData.AllDay,
		eventData.Reminder,
		eventData.ReminderTime,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Start,
		&endTime,
		&event.AllDay,
		&event.Reminder,
		&reminderTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}

	if endTime.Valid {
		event.End = &endTime.Time
	}
	if reminderTime.Valid {
		event.ReminderTime = &reminderTime.Time
	}

	return &event, nil
}

// GetEvent returns an event by ID
func (rp *RpPostgres) GetEvent(id int) (*model.Event, error) {
	query := `
		SELECT id, title, description, start, end, all_day, reminder, reminder_time, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event model.Event
	var endTime sql.NullTime
	var reminderTime sql.NullTime

	err := rp.db.DB.QueryRowContext(context.Background(), query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Start,
		&endTime,
		&event.AllDay,
		&event.Reminder,
		&reminderTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found: %d", id)
		}
		return nil, fmt.Errorf("get event: %w", err)
	}

	if endTime.Valid {
		event.End = &endTime.Time
	}
	if reminderTime.Valid {
		event.ReminderTime = &reminderTime.Time
	}

	return &event, nil
}

// UpdateEvent updates an existing event
func (rp *RpPostgres) UpdateEvent(id int, eventData *model.CreateEventRequest) (*model.Event, error) {
	query := `
		UPDATE events 
		SET title = $2, description = $3, start = $4, end = $5, all_day = $6, reminder = $7, reminder_time = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, title, description, start, end, all_day, reminder, reminder_time, created_at, updated_at
	`

	var event model.Event
	var endTime sql.NullTime
	var reminderTime sql.NullTime

	err := rp.db.DB.QueryRowContext(context.Background(),
		query,
		id,
		eventData.Title,
		eventData.Description,
		eventData.Start,
		eventData.End,
		eventData.AllDay,
		eventData.Reminder,
		eventData.ReminderTime,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Start,
		&endTime,
		&event.AllDay,
		&event.Reminder,
		&reminderTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found: %d", id)
		}
		return nil, fmt.Errorf("update event: %w", err)
	}

	if endTime.Valid {
		event.End = &endTime.Time
	}
	if reminderTime.Valid {
		event.ReminderTime = &reminderTime.Time
	}

	return &event, nil
}

// DeleteEvent deletes an event by ID
func (rp *RpPostgres) DeleteEvent(id int) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := rp.db.DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found: %d", id)
	}

	return nil
}

// GetDayEvents returns events for a specific day
func (rp *RpPostgres) GetDayEvents(date string) ([]model.Event, error) {
	query := `
		SELECT id, title, description, start, end, all_day, reminder, reminder_time, created_at, updated_at
		FROM events
		WHERE DATE(start) = $1
		ORDER BY start ASC
	`

	rows, err := rp.db.DB.QueryContext(context.Background(), query, date)
	if err != nil {
		return nil, fmt.Errorf("query day events: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		var reminderTime sql.NullTime
		var endTime sql.NullTime

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Start,
			&endTime,
			&event.AllDay,
			&event.Reminder,
			&reminderTime,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		if endTime.Valid {
			event.End = &endTime.Time
		}
		if reminderTime.Valid {
			event.ReminderTime = &reminderTime.Time
		}

		events = append(events, event)
	}

	return events, nil
}
