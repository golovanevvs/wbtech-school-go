package customerrors

import "errors"

var (
	ErrContentTypeAJ   = errors.New("content type must be application/json")
	ErrEmptyUserID     = errors.New("user_id must not be empty")
	ErrEmptyID         = errors.New("id must not be empty")
	ErrEmptyTitle      = errors.New("title must not be empty")
	ErrEmptyDate       = errors.New("date must not be empty")
	ErrUserIDNotFound  = errors.New("user id not found")
	ErrEventIDNotFound = errors.New("event id not found")
)
