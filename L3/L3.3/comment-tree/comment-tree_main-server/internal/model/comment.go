package model

import "time"

type Comment struct {
	ID        int
	ParentID  *int
	Text      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Children  []*Comment
}
