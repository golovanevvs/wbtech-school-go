package model

import (
	"time"
)

// Event represents a calendar event
type Event struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Description  string     `json:"description" db:"description"`
	Start        time.Time  `json:"start" db:"start"`
	End          *time.Time `json:"end,omitempty" db:"end"`
	AllDay       bool       `json:"allDay" db:"all_day"`
	Reminder     bool       `json:"reminder" db:"reminder"`
	ReminderTime *time.Time `json:"reminderTime,omitempty" db:"reminder_time"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	Start        time.Time  `json:"start"`
	End          *time.Time `json:"end,omitempty"`
	AllDay       bool       `json:"allDay"`
	Reminder     bool       `json:"reminder"`
	ReminderTime *time.Time `json:"reminderTime,omitempty"`
}

// MonthEventsResponse represents the response for getting events for a month
type MonthEventsResponse struct {
	Events []Event `json:"events"`
}

// CreateEventResponse represents the response for creating an event
type CreateEventResponse struct {
	Event Event `json:"event"`
}

// NotificationChannels represents notification delivery channels
type NotificationChannels struct {
	Telegram bool
	Email    bool
}

// Notice represents a notification to be sent
type Notice struct {
	UserID   int
	Message  string
	Channels NotificationChannels
}
