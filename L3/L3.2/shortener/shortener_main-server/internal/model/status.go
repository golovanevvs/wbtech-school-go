package model

import (
	"errors"
)

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusDeleted   Status = "deleted"
)

func (s Status) Validate() error {
	switch s {
	case StatusScheduled, StatusPending, StatusSent, StatusFailed, StatusDeleted:
	default:
		return errors.New("invalid status: " + string(s))
	}
	return nil
}
