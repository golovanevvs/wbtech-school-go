package model

import "time"

type ShortURL struct {
	ID        int
	Original  string
	Short     string
	Custom    bool
	CreatedAt time.Time
}
