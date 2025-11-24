package eventHandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// ISvForEventHandler interface for event handler service
type ISvForEventHandler interface {
	Create(ctx context.Context, event *model.Event) (*model.Event, error)
	GetByID(ctx context.Context, id int) (*model.Event, error)
	GetAll(ctx context.Context) ([]*model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, id int) error
	UpdateAvailablePlaces(ctx context.Context, eventID int, newAvailablePlaces int) error
	GetByOwnerID(ctx context.Context, ownerID int) ([]*model.Event, error)
}

// EventHandler handles event requests
type EventHandler struct {
	lg *zlog.Zerolog
	sv ISvForEventHandler
}

// New creates a new EventHandler
func New(parentLg *zlog.Zerolog, sv ISvForEventHandler) *EventHandler {
	lg := parentLg.With().Str("component", "handler-eventHandler").Logger()
	return &EventHandler{
		lg: &lg,
		sv: sv,
	}
}

// RegisterPublicRoutes registers the public event routes (no authentication required)
func (hd *EventHandler) RegisterPublicRoutes(rt *ginext.RouterGroup) {
	events := rt.Group("/events")
	{
		events.GET("", hd.GetAll)      // Публичный доступ - список мероприятий
		events.GET("/:id", hd.GetByID) // Публичный доступ - детали мероприятия
	}
}

// RegisterProtectedRoutes registers the protected event routes (authentication required)
func (hd *EventHandler) RegisterProtectedRoutes(rt *ginext.RouterGroup) {
	events := rt.Group("/events")
	{
		events.POST("", hd.Create)                                    // Только для авторизованных - создание мероприятия
		events.PUT("/:id", hd.Update)                                 // Только для авторизованных - обновление мероприятия
		events.DELETE("/:id", hd.Delete)                              // Только для авторизованных - удаление мероприятия
		events.PUT("/:id/available-places", hd.UpdateAvailablePlaces) // Только для авторизованных - обновление доступных мест
	}
}

// Create handles event creation
func (hd *EventHandler) Create(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Create").Logger()

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	event.OwnerID = userIDInt
	event.AvailablePlaces = event.TotalPlaces

	createdEvent, err := hd.sv.Create(c.Request.Context(), &event)
	if err != nil {
		lg.Warn().Err(err).Str("title", event.Title).Msgf("%s failed to create event", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("event_id", createdEvent.ID).Int("owner_id", createdEvent.OwnerID).Str("title", createdEvent.Title).Msgf("%s event created successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusCreated, createdEvent)
}

// GetByID handles getting an event by ID
func (hd *EventHandler) GetByID(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetByID").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid event ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	event, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get event", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("id", id).Msgf("%s event retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, event)
}

// GetAll handles getting all events
func (hd *EventHandler) GetAll(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetAll").Logger()

	events, err := hd.sv.GetAll(c.Request.Context())
	if err != nil {
		lg.Warn().Err(err).Msgf("%s failed to get all events", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("count", len(events)).Msgf("%s events retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, events)
}

// Update handles event update
func (hd *EventHandler) Update(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Update").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid event ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	existingEvent, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get event for ownership check", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	if existingEvent.OwnerID != userIDInt {
		lg.Warn().Int("user_id", userIDInt).Int("event_owner_id", existingEvent.OwnerID).Msgf("%s user is not the owner of the event", pkgConst.Warn)
		c.JSON(http.StatusForbidden, ginext.H{"error": "You are not the owner of this event"})
		return
	}

	event.ID = id

	err = hd.sv.Update(c.Request.Context(), &event)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Str("title", event.Title).Msgf("%s failed to update event", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	updatedEvent, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get updated event", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to retrieve updated event"})
		return
	}

	lg.Debug().Int("id", id).Str("title", event.Title).Msgf("%s event updated successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, updatedEvent)
}

// Delete handles event deletion
func (hd *EventHandler) Delete(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Delete").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid event ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	existingEvent, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get event for ownership check", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	if existingEvent.OwnerID != userIDInt {
		lg.Warn().Int("user_id", userIDInt).Int("event_owner_id", existingEvent.OwnerID).Msgf("%s user is not the owner of the event", pkgConst.Warn)
		c.JSON(http.StatusForbidden, ginext.H{"error": "You are not the owner of this event"})
		return
	}

	err = hd.sv.Delete(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to delete event", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("id", id).Msgf("%s event deleted successfully", pkgConst.OpSuccess)
	c.Status(http.StatusOK)
}

// UpdateAvailablePlaces handles updating available places for an event
func (hd *EventHandler) UpdateAvailablePlaces(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "UpdateAvailablePlaces").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid event ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	var req struct {
		NewAvailablePlaces int `json:"new_available_places" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	existingEvent, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get event for ownership check", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	if existingEvent.OwnerID != userIDInt {
		lg.Warn().Int("user_id", userIDInt).Int("event_owner_id", existingEvent.OwnerID).Msgf("%s user is not the owner of the event", pkgConst.Warn)
		c.JSON(http.StatusForbidden, ginext.H{"error": "You are not the owner of this event"})
		return
	}

	err = hd.sv.UpdateAvailablePlaces(c.Request.Context(), id, req.NewAvailablePlaces)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Int("new_available_places", req.NewAvailablePlaces).Msgf("%s failed to update available places", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("id", id).Int("new_available_places", req.NewAvailablePlaces).Msgf("%s available places updated successfully", pkgConst.OpSuccess)
	c.Status(http.StatusOK)
}
