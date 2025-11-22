package model

import (
	"time"
)

type User struct {
	ID                    int       `json:"id" db:"id"`
	Email                 string    `json:"email" db:"email"`
	Name                  string    `json:"name" db:"name"`
	PasswordHash          string    `json:"-" db:"password_hash"`
	TelegramUsername      *string   `json:"telegramUsername,omitempty" db:"telegram_username"`
	TelegramChatID        *int64    `json:"telegramChatID,omitempty" db:"telegram_chat_id"`
	TelegramNotifications bool      `json:"telegramNotifications" db:"telegram_notifications"`
	EmailNotifications    bool      `json:"emailNotifications" db:"email_notifications"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

type Event struct {
	ID                    int       `json:"id" db:"id"`
	Title                 string    `json:"title" db:"title"`
	Date                  time.Time `json:"date" db:"date"`
	Description           string    `json:"description" db:"description"`
	TotalPlaces           int       `json:"total_places" db:"total_places"`
	AvailablePlaces       int       `json:"available_places" db:"available_places"`
	BookingDeadline       int       `json:"booking_deadline" db:"booking_deadline"` // minute
	OwnerID               int       `json:"owner_id" db:"owner_id"`
	TelegramNotifications bool      `json:"telegram_notifications" db:"telegram_notifications"`
	EmailNotifications    bool      `json:"email_notifications" db:"email_notifications"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
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
	ConfirmedAt *time.Time    `json:"confirmedAt,omitempty" db:"confirmed_at"`
	CancelledAt *time.Time    `json:"cancelledAt,omitempty" db:"cancelled_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type NotificationChannels struct {
	Telegram bool
	Email    bool
}

type Notice struct {
	UserID   int
	Message  string
	Channels NotificationChannels
}
