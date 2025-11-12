package model

import "time"

type ShortURL struct {
	ID        int
	Original  string `validate:"required,url"`
	Short     string
	Custom    bool
	CreatedAt time.Time
}
