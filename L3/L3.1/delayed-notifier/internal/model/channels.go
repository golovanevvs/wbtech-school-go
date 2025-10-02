package model

import (
	"encoding/json"
	"errors"
)

type Channels []ChannelInfo

type ChannelInfo struct {
	Type  ChannelType `json:"type"`
	Value string      `json:"value"`
}

type ChannelType string

const (
	ChannelEmail    ChannelType = "email"
	ChannelTelegram ChannelType = "telegram"
)

func (c Channels) Validate() error {
	if len(c) == 0 {
		return errors.New("channels cannot be empty")
	}

	for _, ch := range c {
		switch ch.Type {
		case ChannelEmail, ChannelTelegram:
		default:
			return errors.New("invalid channel type: " + string(ch.Type))
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
		switch ch.Type {
		case ChannelEmail, ChannelTelegram:
		default:
			return errors.New("invalid channel type: " + string(ch.Type))
		}
	}

	*c = Channels(s)

	return nil
}
