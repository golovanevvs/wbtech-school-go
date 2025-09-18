package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/customerrors"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repository struct {
	store  map[string]map[string]model.Event
	mu     sync.RWMutex
	logger zerolog.Logger
}

func New(logger *zerolog.Logger) *Repository {
	return &Repository{
		store:  make(map[string]map[string]model.Event),
		mu:     sync.RWMutex{},
		logger: logger.With().Str("component", "repository").Logger(),
	}
}

func (r *Repository) Create(event model.Event) (id string) {
	log := r.logger.With().Str("method", "Create").Str("user_id", event.UserId).Logger()
	id = uuid.New().String()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[event.UserId]; !ok {
		r.store[event.UserId] = make(map[string]model.Event)
	}

	r.store[event.UserId][id] = event
	log.Info().Str("event_id", id).Msg("event created")

	return
}

func (r *Repository) Update(event model.Event) error {
	log := r.logger.With().Str("method", "Update").Str("user_id", event.UserId).Logger()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[event.UserId]; !ok {
		log.Error().Msg(customerrors.ErrUserIDNotFound.Error())
		return customerrors.ErrUserIDNotFound
	}

	if _, ok := r.store[event.UserId][event.Id]; !ok {
		log.Error().Msg(customerrors.ErrEventIDNotFound.Error())
		return customerrors.ErrEventIDNotFound
	}

	r.store[event.UserId][event.Id] = event
	log.Info().Str("event_id", event.Id).Msg("event updated")

	return nil
}

func (r *Repository) Delete(event model.Event) error {
	log := r.logger.With().Str("method", "Delete").Str("user_id", event.UserId).Logger()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[event.UserId]; !ok {
		log.Error().Msg(customerrors.ErrUserIDNotFound.Error())
		return customerrors.ErrUserIDNotFound
	}

	if _, ok := r.store[event.UserId][event.Id]; !ok {
		log.Error().Msg(customerrors.ErrEventIDNotFound.Error())
		return customerrors.ErrEventIDNotFound
	}

	delete(r.store[event.UserId], event.Id)
	log.Info().Str("event_id", event.Id).Msg("event deleted")

	return nil
}

func (r *Repository) LoadForDay(userId string, day time.Time) ([]model.Event, error) {
	log := r.logger.With().Str("method", "GetForDay").Str("user_id", userId).Logger()
	events := make([]model.Event, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[userId]; !ok {
		log.Error().Msg("failed user id")
		return events, fmt.Errorf("error get events for day: failed user id")
	}

	for _, v := range r.store[userId] {
		y1, m1, d1 := time.Time(v.Date).Date()
		y2, m2, d2 := day.Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			events = append(events, v)
		}
	}

	log.Info().Int("events_count", len(events)).Msg("events retrieved for day")
	return events, nil
}

func (r *Repository) LoadForWeek(userId string, week time.Time) ([]model.Event, error) {
	log := r.logger.With().Str("method", "LoadForWeek").Str("user_id", userId).Logger()
	events := make([]model.Event, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[userId]; !ok {
		log.Error().Msg(customerrors.ErrUserIDNotFound.Error())
		return events, customerrors.ErrUserIDNotFound
	}

	startOfWeek := week.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	for _, v := range r.store[userId] {
		if time.Time(v.Date).After(startOfWeek) && time.Time(v.Date).Before(endOfWeek) {
			events = append(events, v)
		}
	}

	log.Info().Int("events_count", len(events)).Msg("events retrieved for week")
	return events, nil
}

func (r *Repository) LoadForMonth(userId string, month time.Time) ([]model.Event, error) {
	log := r.logger.With().Str("method", "LoadForMonth").Str("user_id", userId).Logger()
	events := make([]model.Event, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[userId]; !ok {
		log.Error().Msg(customerrors.ErrUserIDNotFound.Error())
		return events, customerrors.ErrUserIDNotFound
	}

	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	for _, v := range r.store[userId] {
		if time.Time(v.Date).After(startOfMonth) && time.Time(v.Date).Before(endOfMonth) {
			events = append(events, v)
		}
	}

	log.Info().Int("events_count", len(events)).Msg("events retrieved for month")
	return events, nil
}
