package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Notice struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id" validate:"required"`
	Message   string     `json:"message" validate:"required"`
	Channels  Channels   `json:"channels"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty" validate:"omitempty,gtfield=CreatedAt"`
	Status    Status     `json:"status"`
}

func (n *Notice) Validate() error {
	validate := validator.New()

	if err := validate.Struct(n); err != nil {
		return err
	}

	if err := n.Channels.Validate(); err != nil {
		return err
	}

	if err := n.Status.Validate(); err != nil {
		return err
	}

	return nil
}
