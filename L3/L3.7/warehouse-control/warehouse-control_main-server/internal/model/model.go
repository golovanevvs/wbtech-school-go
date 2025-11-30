package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	UserName     string    `json:"username" db:"username"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	UserRole     string    `json:"user_role" db:"user_role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Item struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ItemAction struct {
	ID         int       `json:"id" db:"id"`
	ItemID     int       `json:"item_id" db:"item_id"`
	ActionType string    `json:"action_type" db:"action_type"`
	UserID     int       `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	Changes    string    `json:"changes" db:"changes"` // JSONB в PostgreSQL
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
