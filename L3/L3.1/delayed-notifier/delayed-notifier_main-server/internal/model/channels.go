package model

import (
	"encoding/json"
	"errors"
)

var (
	ErrEmptyChannels   = errors.New("channels cannot be empty")
	ErrInvalidTypeBase = errors.New("invalid channel type")
	ErrEmptyValue      = errors.New("channel value cannot be empty")
)

type Channels []ChannelInfo

type ChannelInfo struct {
	Type  ChannelType `json:"type" validate:"required,oneof=email telegram" binding:"required,oneof=email telegram"`
	Value string      `json:"value" validate:"required" binding:"required"`
}

type ChannelType string

const (
	ChannelEmail    ChannelType = "email"
	ChannelTelegram ChannelType = "telegram"
)

func (t ChannelType) IsValid() bool {
	switch t {
	case ChannelEmail, ChannelTelegram:
		return true
	}
	return false
}

func (c Channels) Validate() error {
	if len(c) == 0 {
		return ErrEmptyChannels
	}

	validCount := 0
	for _, ch := range c {
		if ch.Type.IsValid() && ch.Value != "" {
			validCount++
		}
	}

	if validCount == 0 {
		return ErrEmptyValue
	}

	return nil
}

func (c *Channels) UnmarshalJSON(data []byte) error {
	var s []ChannelInfo
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ch := Channels(s)
	if err := ch.Validate(); err != nil {
		return err
	}

	*c = ch
	return nil
}
