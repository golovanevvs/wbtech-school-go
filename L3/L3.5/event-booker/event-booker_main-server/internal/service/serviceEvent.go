package service

import (
	"context"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
)

// EventRepo interface for event repository
type EventRepo interface {
	Create(event *model.Event) (*model.Event, error)
	GetByID(id int) (*model.Event, error)
	GetAll() ([]*model.Event, error)
	Update(event *model.Event) error
	Delete(id int) error
	UpdateAvailablePlaces(eventID int, newAvailablePlaces int) error
}

// EventService service for working with events
type EventService struct {
	rp EventRepo
}

// NewEventService creates a new EventService
func NewEventService(rp EventRepo) *EventService {
	return &EventService{rp: rp}
}

// Create creates a new event
func (s *EventService) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	// Set creation and update timestamps
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	// Ensure available places equal total places initially
	event.AvailablePlaces = event.TotalPlaces

	// Create the event in the repository
	createdEvent, err := s.rp.Create(event)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// GetByID returns an event by ID
func (s *EventService) GetByID(ctx context.Context, id int) (*model.Event, error) {
	return s.rp.GetByID(id)
}

// GetAll returns all events
func (s *EventService) GetAll(ctx context.Context) ([]*model.Event, error) {
	return s.rp.GetAll()
}

// Update updates an event
func (s *EventService) Update(ctx context.Context, event *model.Event) error {
	// Update the timestamp
	event.UpdatedAt = time.Now()

	return s.rp.Update(event)
}

// Delete deletes an event by ID
func (s *EventService) Delete(ctx context.Context, id int) error {
	return s.rp.Delete(id)
}

// UpdateAvailablePlaces updates the available places for an event
func (s *EventService) UpdateAvailablePlaces(ctx context.Context, eventID int, newAvailablePlaces int) error {
	return s.rp.UpdateAvailablePlaces(eventID, newAvailablePlaces)
}
