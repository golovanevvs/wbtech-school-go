package pkgErrors

import (
	"errors"
)

var (
	ErrContentTypeAJ  = errors.New("content type must be application/json")
	ErrEmptyUserID    = errors.New("user_id must not be empty")
	ErrEmptyID        = errors.New("id must not be empty")
	ErrEmptyTitle     = errors.New("title must not be empty")
	ErrEmptyDate      = errors.New("date must not be empty")
	ErrUserNotFound   = errors.New("user id not found")
	ErrNoticeNotFound = errors.New("notice not found")
)
