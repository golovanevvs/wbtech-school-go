package eventHandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// IService interface for event handler service
type IService interface {
	GetMonthEvents(ctx context.Context, year, month int) ([]model.Event, error)
	CreateEvent(ctx context.Context, eventData *model.CreateEventRequest) (*model.Event, error)
	GetEvent(ctx context.Context, id int) (*model.Event, error)
	UpdateEvent(ctx context.Context, id int, eventData *model.CreateEventRequest) (*model.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	GetDayEvents(ctx context.Context, date string) ([]model.Event, error)
}

// EventHandler handles calendar event HTTP requests
type EventHandler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

// NewEventHandler creates a new EventHandler
func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *EventHandler {
	lg := parentLg.With().Str("component", "eventHandler").Logger()
	return &EventHandler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

// RegisterRoutes registers the event handler routes
func (hd *EventHandler) RegisterRoutes() {
	events := hd.rt.Group("/events")
	{
		events.GET("/month", hd.GetMonthEventsHandler)
		events.GET("/day", hd.GetDayEventsHandler)
	}

	event := hd.rt.Group("/event")
	{
		event.POST("/create", hd.CreateEventHandler)
		event.GET("/:id", hd.GetEventHandler)
		event.PUT("/:id", hd.UpdateEventHandler)
		event.DELETE("/:id", hd.DeleteEventHandler)
	}
}

// GetMonthEventsHandler gets events for a specific month
func (hd *EventHandler) GetMonthEventsHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetMonthEventsHandler").Logger()

	// Parse year and month from query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		lg.Warn().Str("year", yearStr).Str("month", monthStr).Msg("Missing year or month parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month parameters are required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		lg.Warn().Str("year", yearStr).Err(err).Msg("Invalid year parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		lg.Warn().Str("month", monthStr).Err(err).Msg("Invalid month parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month parameter"})
		return
	}

	if month < 1 || month > 12 {
		lg.Warn().Int("month", month).Msg("Month out of range")
		c.JSON(http.StatusBadRequest, gin.H{"error": "month must be between 1 and 12"})
		return
	}

	// Get events from service
	events, err := hd.sv.GetMonthEvents(c.Request.Context(), year, month)
	if err != nil {
		lg.Error().Err(err).Int("year", year).Int("month", month).Msg("Failed to get month events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get events"})
		return
	}

	lg.Debug().Int("year", year).Int("month", month).Int("events_count", len(events)).Msg("Successfully retrieved month events")

	// Return response in the format expected by the client
	response := model.MonthEventsResponse{
		Events: events,
	}
	c.JSON(http.StatusOK, response)
}

// CreateEventHandler creates a new event
func (hd *EventHandler) CreateEventHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "CreateEventHandler").Logger()

	var eventData model.CreateEventRequest
	if err := c.ShouldBindJSON(&eventData); err != nil {
		lg.Warn().Err(err).Msg("Invalid JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload: " + err.Error()})
		return
	}

	// Validate required fields
	if eventData.Title == "" {
		lg.Warn().Msg("Title is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if eventData.Start.IsZero() {
		lg.Warn().Msg("Start time is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time is required"})
		return
	}

	// If end time is provided, it should be after start time
	if eventData.End != nil && !eventData.End.IsZero() && eventData.End.Before(eventData.Start) {
		lg.Warn().Time("start", eventData.Start).Time("end", *eventData.End).Msg("End time is before start time")
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time must be after start time"})
		return
	}

	// Create event through service
	event, err := hd.sv.CreateEvent(c.Request.Context(), &eventData)
	if err != nil {
		lg.Error().Err(err).Str("title", eventData.Title).Msg("Failed to create event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	lg.Info().Int("event_id", event.ID).Str("title", event.Title).Msg("Event created successfully")

	response := model.CreateEventResponse{
		Event: *event,
	}
	c.JSON(http.StatusCreated, response)
}

// GetEventHandler gets a single event by ID
func (hd *EventHandler) GetEventHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetEventHandler").Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		lg.Warn().Str("id", idStr).Err(err).Msg("Invalid event ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := hd.sv.GetEvent(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Int("event_id", id).Err(err).Msg("Event not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEventHandler updates an existing event
func (hd *EventHandler) UpdateEventHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "UpdateEventHandler").Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		lg.Warn().Str("id", idStr).Err(err).Msg("Invalid event ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var eventData model.CreateEventRequest
	if err := c.ShouldBindJSON(&eventData); err != nil {
		lg.Warn().Err(err).Msg("Invalid JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload: " + err.Error()})
		return
	}

	event, err := hd.sv.UpdateEvent(c.Request.Context(), id, &eventData)
	if err != nil {
		lg.Error().Err(err).Int("event_id", id).Msg("Failed to update event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	lg.Info().Int("event_id", event.ID).Str("title", event.Title).Msg("Event updated successfully")

	c.JSON(http.StatusOK, event)
}

// DeleteEventHandler deletes an event
func (hd *EventHandler) DeleteEventHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "DeleteEventHandler").Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		lg.Warn().Str("id", idStr).Err(err).Msg("Invalid event ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	err = hd.sv.DeleteEvent(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Int("event_id", id).Err(err).Msg("Event not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	lg.Info().Int("event_id", id).Msg("Event deleted successfully")

	c.Status(http.StatusNoContent)
}

// GetDayEventsHandler gets events for a specific day
func (hd *EventHandler) GetDayEventsHandler(c *ginext.Context) {
	lg := hd.lg.With().Str("method", "GetDayEventsHandler").Logger()

	dateStr := c.Query("date")
	if dateStr == "" {
		lg.Warn().Msg("Missing date parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	// Validate date format (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		lg.Warn().Str("date", dateStr).Err(err).Msg("Invalid date format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in YYYY-MM-DD format"})
		return
	}

	events, err := hd.sv.GetDayEvents(c.Request.Context(), dateStr)
	if err != nil {
		lg.Error().Err(err).Str("date", dateStr).Msg("Failed to get day events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get events"})
		return
	}

	lg.Debug().Str("date", dateStr).Int("events_count", len(events)).Msg("Successfully retrieved day events")

	c.JSON(http.StatusOK, gin.H{"events": events})
}
