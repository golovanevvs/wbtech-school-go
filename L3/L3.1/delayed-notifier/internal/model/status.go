package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type StatusByChannel map[ChannelInfo]StatusType

type StatusType string

const (
	StatusScheduled StatusType = "scheduled"
	StatusPending   StatusType = "pending"
	StatusSent      StatusType = "sent"
	StatusFailed    StatusType = "failed"
)

func (s StatusByChannel) Validate() error {
	for ch, status := range s {
		switch ch.Type {
		case ChannelEmail, ChannelTelegram:
		default:
			return errors.New("invalid channel type: " + string(ch.Type))
		}

		if ch.Value == "" {
			return errors.New("channel value cannot be empty")
		}

		switch status {
		case StatusScheduled, StatusPending, StatusSent, StatusFailed:
		default:
			return errors.New("invalid status: " + string(status))
		}
	}

	return nil
}

func (s StatusByChannel) MarshalJSON() ([]byte, error) {
	raw := make(map[string]string)
	for ch, status := range s {
		key := fmt.Sprintf("%s:%s", ch.Type, ch.Value)
		raw[key] = string(status)
	}

	return json.Marshal(raw)
}
