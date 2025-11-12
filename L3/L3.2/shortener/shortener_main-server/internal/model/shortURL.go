package model

import "time"

type ShortURL struct {
	ID        int       `json:"id"`
	Original  string    `json:"original" validate:"required,url"`
	Short     string    `json:"short"`
	Custom    bool      `json:"custom"`
	CreatedAt time.Time `json:"created_at"`
}
