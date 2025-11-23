package bookingHandler

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

// ISvForBookingHandler interface for booking handler service
type ISvForBookingHandler interface {
	Create(ctx context.Context, userID int, eventID int, bookingDeadlineMinutes int) (*model.Booking, error)
	GetByID(ctx context.Context, id int) (*model.Booking, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.Booking, error)
	GetByUserIDAndEventID(ctx context.Context, userID int, eventID int) (*model.Booking, error)
	Confirm(ctx context.Context, bookingID int) error
	Cancel(ctx context.Context, bookingID int) error
}

// BookingHandler handles booking requests
type BookingHandler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv ISvForBookingHandler
}

// NewBookingHandler creates a new BookingHandler
func New(parentLg *zlog.Zerolog, sv ISvForBookingHandler) *BookingHandler {
	lg := parentLg.With().Str("component", "handler-bookingHandler").Logger()
	return &BookingHandler{
		lg: &lg,
		sv: sv,
	}
}

// RegisterRoutes registers the booking routes
func (hd *BookingHandler) RegisterRoutes(rt *ginext.RouterGroup) {
	bookings := rt.Group("/bookings")
	{
		bookings.POST("", hd.Create)
		bookings.GET("", hd.GetCurrentUserBookings)
		bookings.GET("/:id", hd.GetByID)
		bookings.GET("/user/:userID", hd.GetByUserID)
		bookings.GET("/user/event/:eventID", hd.GetByUserIDAndEventID)
		bookings.POST("/:id/confirm", hd.Confirm)
		bookings.POST("/:id/cancel", hd.Cancel)
	}
}

// Create handles booking creation
func (hd *BookingHandler) Create(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Create").Logger()

	var req struct {
		EventID                int `json:"event_id" binding:"required"`
		BookingDeadlineMinutes int `json:"booking_deadline_minutes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	// Получаем user_id из контекста как int
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	booking, err := hd.sv.Create(c.Request.Context(), userID, req.EventID, req.BookingDeadlineMinutes)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userID).Int("event_id", req.EventID).Msgf("%s failed to create booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("booking_id", booking.ID).Int("user_id", userID).Int("event_id", req.EventID).Msgf("%s booking created successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusCreated, booking)
}

// GetByID handles getting a booking by ID
func (hd *BookingHandler) GetByID(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetByID").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid booking ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	booking, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("id", id).Msgf("%s booking retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, booking)
}

// GetByUserID handles getting bookings by user ID
func (hd *BookingHandler) GetByUserID(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetByUserID").Logger()

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		lg.Warn().Err(err).Str("user_id", c.Param("userID")).Msgf("%s invalid user ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	bookings, err := hd.sv.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userID).Msgf("%s failed to get bookings", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("user_id", userID).Int("count", len(bookings)).Msgf("%s bookings retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, bookings)
}

// GetByUserIDAndEventID handles getting a booking by user ID and event ID
func (hd *BookingHandler) GetByUserIDAndEventID(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetByUserIDAndEventID").Logger()

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		lg.Warn().Err(err).Str("event_id", c.Param("eventID")).Msgf("%s invalid event ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	booking, err := hd.sv.GetByUserIDAndEventID(c.Request.Context(), userID, eventID)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userID).Int("event_id", eventID).Msgf("%s failed to get booking", pkgConst.Warn)
		c.JSON(http.StatusNotFound, ginext.H{"error": "Booking not found"})
		return
	}

	lg.Debug().Int("booking_id", booking.ID).Int("user_id", userID).Int("event_id", eventID).Msgf("%s booking retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, booking)
}

// Confirm handles booking confirmation
func (hd *BookingHandler) Confirm(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Confirm").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid booking ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	err = hd.sv.Confirm(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to confirm booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	// Получаем обновленную бронь для возврата клиенту
	updatedBooking, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get updated booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to retrieve updated booking"})
		return
	}

	lg.Debug().Int("id", id).Msgf("%s booking confirmed successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, updatedBooking)
}

// Cancel handles booking cancellation
func (hd *BookingHandler) Cancel(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Cancel").Logger()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		lg.Warn().Err(err).Str("id", c.Param("id")).Msgf("%s invalid booking ID", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrInvalidID.Error()})
		return
	}

	err = hd.sv.Cancel(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to cancel booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	// Получаем обновленную бронь для возврата клиенту
	updatedBooking, err := hd.sv.GetByID(c.Request.Context(), id)
	if err != nil {
		lg.Warn().Err(err).Int("id", id).Msgf("%s failed to get updated booking", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to retrieve updated booking"})
		return
	}

	lg.Debug().Int("id", id).Msgf("%s booking cancelled successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, updatedBooking)
}

// GetCurrentUserBookings handles getting bookings for the current user
func (hd *BookingHandler) GetCurrentUserBookings(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetCurrentUserBookings").Logger()

	// Получаем user_id из контекста
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		lg.Warn().Msg("User ID is not of type int")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	bookings, err := hd.sv.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userID).Msgf("%s failed to get bookings", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("user_id", userID).Int("count", len(bookings)).Msgf("%s bookings retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, bookings)
}
