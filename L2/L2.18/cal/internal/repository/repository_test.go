package repository

import (
	"testing"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
	"github.com/stretchr/testify/assert"
)

type expected struct {
	id    string
	error error
}

func TestCreate(t *testing.T) {
	zlog.Init()
	rp := New(&zlog.Logger)

	date, _ := time.Parse("2006-01-02", "2025-09-18")

	id := rp.Create(model.Event{
		UserId:  "user-id-001",
		Title:   "Title1",
		Comment: "Comment1",
		Date:    model.DateOnly(date),
	})

	assert.NotEmpty(t, id)
}

func TestUpdate(t *testing.T) {
	zlog.Init()
	rp := New(&zlog.Logger)

	tests := []struct {
		name     string
		input    model.Event
		expected expected
	}{

		{
			name: "test1",
			input: model.Event{
				UserId:  "user-id-001",
				Title:   "Title1",
				Comment: "Comment1",
			},
			expected: expected{
				id:    "sdfs",
				error: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "test1":
				date, _ := time.Parse("2006-01-02", "2025-09-18")
				test.input.Date = model.DateOnly(date)

				id := rp.Create(test.input)

				assert.NotEmpty(t, id)
			}
		})
	}
}
