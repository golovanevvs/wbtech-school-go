package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Event struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Date            time.Time `json:"date" db:"date"`
	Description     string    `json:"description" db:"description"`
	TotalPlaces     int       `json:"total_places" db:"total_places"`
	AvailablePlaces int       `json:"available_places" db:"available_places"`
	BookingDeadline int       `json:"booking_deadline" db:"booking_deadline"` // minute
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          int           `json:"id" db:"id"`
	UserID      int           `json:"user_id" db:"user_id"`
	EventID     int           `json:"event_id" db:"event_id"`
	Status      BookingStatus `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time     `json:"expires_at" db:"expires_at"`
	ConfirmedAt *time.Time    `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CancelledAt *time.Time    `json:"cancelled_at,omitempty" db:"cancelled_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
