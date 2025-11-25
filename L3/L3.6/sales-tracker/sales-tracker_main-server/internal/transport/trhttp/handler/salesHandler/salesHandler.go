package salesHandler

import (
	"context"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// IService defines the service interface for SalesRecord operations
type IService interface {
	CreateSalesRecord(ctx context.Context, data model.Data) (int, error)
}

// SalesHandler handles HTTP requests for SalesRecord
type SalesHandler struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

// New creates a new SalesHandler instance
func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *SalesHandler {
	lg := parentLg.With().Str("component", "SalesHandler").Logger()
	return &SalesHandler{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

// RegisterRoutes registers routes for SalesRecord operations
func (h *SalesHandler) RegisterRoutes() {
	h.rt.POST("/items", h.CreateSalesRecord)
}

// CreateSalesRecord handles POST /items for creating a new record
func (h *SalesHandler) CreateSalesRecord(c *ginext.Context) {
	lg := h.lg.With().Str("method", "CreateSalesRecord").Logger()

	var request CreateSalesRecordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		lg.Warn().Err(err).Msgf("%s failed to bind JSON", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	if request.Type == "" {
		lg.Warn().Msgf("%s Type is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Type is required",
		})
		return
	}

	if request.Category == "" {
		lg.Warn().Msgf("%s Category is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Category is required",
		})
		return
	}

	if request.Date == "" {
		lg.Warn().Msgf("%s Date is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Date is required",
		})
		return
	}

	if request.Amount <= 0 {
		lg.Warn().Float64("amount", request.Amount).Msgf("%s Amount must be positive", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Amount must be positive",
		})
		return
	}

	if request.Type != "income" && request.Type != "expense" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Type must be 'income' or 'expense'",
		})
		return
	}

	serviceRequest := model.Data{
		Type:     request.Type,
		Category: request.Category,
		Date:     request.Date,
		Amount:   request.Amount,
	}

	id, err := h.sv.CreateSalesRecord(c.Request.Context(), serviceRequest)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to create sales record")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create sales record",
		})
		return
	}

	response := CreateSalesRecordResponse{
		ID: id,
	}

	lg.Debug().Int("ID", id).Msgf("%s record added successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusCreated, response)
}
