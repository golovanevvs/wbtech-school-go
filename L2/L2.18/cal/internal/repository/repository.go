package repository

import (
	"sync"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
)

type Repository struct {
	store map[int]model.Event
	mu    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		store: make(map[int]model.Event),
		mu:    sync.RWMutex{},
	}
}
