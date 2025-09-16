package repository

import (
	"sync"

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

func (r *Repository) Create(event model.Event) {
	r.mu.Lock()

	if _, ok := r.store[event.UserId]; !ok {
		r.store[event.UserId] = make(map[string]model.Event)
	}

	id := uuid.New().String()

	r.store[event.UserId][id] = event

	r.mu.Unlock()
}
