package model

import "time"

type Analitic struct {
	ID        int
	ShortID   int
	UserAgent string
	CreatedAt time.Time
}
