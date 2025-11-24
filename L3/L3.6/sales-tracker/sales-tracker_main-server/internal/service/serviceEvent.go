package service

import (
	"context"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
)

// IEventRp interface for event repository
type IEventRp interface {
	Create(event *model.Event) (*model.Event, error)
	GetByID(id int) (*model.Event, error)
	GetAll() ([]*model.Event, error)
	Update(event *model.Event) error
	Delete(id int) error
	UpdateAvailablePlaces(eventID int, newAvailablePlaces int) error
	GetByOwnerID(ownerID int) ([]*model.Event, error)
}

// EventService service for working with events
type EventService struct {
	rp IEventRp
}

// NewEventService creates a new EventService
func NewEventService(rp IEventRp) *EventService {
	return &EventService{rp: rp}
}

// Create creates a new event
func (sv *EventService) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	event.AvailablePlaces = event.TotalPlaces

	createdEvent, err := sv.rp.Create(event)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// GetByID returns an event by ID
func (sv *EventService) GetByID(ctx context.Context, id int) (*model.Event, error) {
	return sv.rp.GetByID(id)
}

// GetAll returns all events
func (sv *EventService) GetAll(ctx context.Context) ([]*model.Event, error) {
	return sv.rp.GetAll()
}

// GetByOwnerID returns all events owned by a user
func (sv *EventService) GetByOwnerID(ctx context.Context, ownerID int) ([]*model.Event, error) {
	return sv.rp.GetByOwnerID(ownerID)
}

// Update updates an event
func (sv *EventService) Update(ctx context.Context, event *model.Event) error {
	// Update the timestamp
	event.UpdatedAt = time.Now()

	return sv.rp.Update(event)
}

// Delete deletes an event by ID
func (sv *EventService) Delete(ctx context.Context, id int) error {
	return sv.rp.Delete(id)
}

// UpdateAvailablePlaces updates the available places for an event
func (sv *EventService) UpdateAvailablePlaces(ctx context.Context, eventID int, newAvailablePlaces int) error {
	return sv.rp.UpdateAvailablePlaces(eventID, newAvailablePlaces)
}
