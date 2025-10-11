package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ReqNotice struct {
	UserID   int        `json:"user_id" binding:"required,numeric"`
	Message  string     `json:"message" binding:"required"`
	Channels Channels   `json:"channels" binding:"required"`
	SentAt   *time.Time `json:"sent_at" validate:"required"`
}

func (n *ReqNotice) Validate() error {
	validate := validator.New()

	if err := validate.Struct(n); err != nil {
		return err
	}

	if err := n.Channels.Validate(); err != nil {
		return err
	}

	return nil
}
