package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Channels []ChannelInfo

type ChannelInfo struct {
	Type  ChannelType `json:"type" validate:"required"`
	Value string      `json:"value" validate:"required"`
}

type ChannelType string

const (
	ChannelEmail    ChannelType = "email"
	ChannelTelegram ChannelType = "telegram"
)

var (
	ErrEmptyChannels   = errors.New("channels cannot be empty")
	ErrInvalidTypeBase = errors.New("invalid channel type")
	ErrEmptyValue      = errors.New("channel value cannot be empty")
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

	for _, ch := range c {
		if !ch.Type.IsValid() {
			return ErrInvalidTypeBase
		}
		if ch.Value == "" {
			return ErrEmptyValue
		}
	}

	return nil
}

func (c *Channels) UnmarshalJSON(data []byte) error {
	var s []ChannelInfo
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	for _, ch := range s {
		if !ch.Type.IsValid() {
			return fmt.Errorf("%w: %s", ErrInvalidTypeBase, ch.Type)
		}
		if ch.Value == "" {
			return ErrEmptyValue
		}
	}

	*c = Channels(s)
	return nil
}
