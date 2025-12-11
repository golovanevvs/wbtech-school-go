package model

import "time"

type Analitic struct {
	ID        int       `json:"id"`
	Short     string    `json:"short"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
