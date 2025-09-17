package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
	"github.com/google/uuid"
)

type Repository struct {
	store map[string]map[string]model.Event
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		store: make(map[string]map[string]model.Event),
		mu:    sync.RWMutex{},
	}
}

func (r *Repository) Create(event model.Event) (id string) {
	id = uuid.New().String()

	r.mu.Lock()

	if _, ok := r.store[event.UserId]; !ok {
		r.store[event.UserId] = make(map[string]model.Event)
	}

	r.store[event.UserId][id] = event

	r.mu.Unlock()

	return
}

func (r *Repository) Update(event model.Event) error {
	r.mu.Lock()

	if _, ok := r.store[event.UserId]; !ok {
		return fmt.Errorf("error update event: failed user id")
	}

	if _, ok := r.store[event.UserId][event.Id]; !ok {
		return fmt.Errorf("error update event: failed event id")
	}

	r.store[event.UserId][event.Id] = event

	r.mu.Unlock()

	return nil
}

func (r *Repository) Delete(event model.Event) error {
	r.mu.Lock()

	if _, ok := r.store[event.UserId]; !ok {
		return fmt.Errorf("error delete event: failed user id")
	}

	if _, ok := r.store[event.UserId][event.Id]; !ok {
		return fmt.Errorf("error delete event: failed event id")
	}

	delete(r.store[event.UserId], event.Id)

	r.mu.Unlock()

	return nil
}

func (r *Repository) GetForDay(userId string, day time.Time) ([]model.Event, error) {
	events := make([]model.Event, 0)

	r.mu.Lock()

	if _, ok := r.store[userId]; !ok {
		return events, fmt.Errorf("error get events for day: failed user id")
	}

	for _, v := range r.store[userId] {
		y1, m1, d1 := time.Time(v.Date).Date()
		y2, m2, d2 := day.Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			events = append(events, v)
		}
	}

	r.mu.Unlock()

	return events, nil
}
