package model

import (
	"encoding/json"
	"time"
)

type dateOnly time.Time

type Event struct {
	UserId  string   `json:"user_id"`
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Comment string   `json:"comment"`
	Date    dateOnly `json:"date"`
}

type Resp struct {
	Id     string `json:"id,omitempty"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (d *dateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = dateOnly(t)
	return nil
}

func (d dateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}
