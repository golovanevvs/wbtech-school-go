package salesHandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// IService defines the service interface for SalesRecord operations
type IService interface {
	CreateSalesRecord(ctx context.Context, data model.Data) (int, error)
	GetSalesRecords(ctx context.Context, sortOptions model.SortOptions) ([]model.Data, error)
	UpdateSalesRecord(ctx context.Context, id int, data model.Data) error
	DeleteSalesRecord(ctx context.Context, id int) error
	GetAnalytics(ctx context.Context, from, to string) (model.Analytics, error)
	ExportCSV(ctx context.Context, from, to string) ([]byte, error)
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
	h.rt.GET("/items", h.GetSalesRecords)
	h.rt.PUT("/items/:id", h.UpdateSalesRecord)
	h.rt.DELETE("/items/:id", h.DeleteSalesRecord)
	h.rt.GET("/analytics", h.GetAnalytics)
	h.rt.GET("/items/export", h.ExportCSV)
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

// GetSalesRecords handles GET /items for retrieving records with sorting
func (h *SalesHandler) GetSalesRecords(c *ginext.Context) {
	lg := h.lg.With().Str("method", "GetSalesRecords").Logger()

	// Get sorting parameters from query
	field := c.DefaultQuery("field", "id")
	direction := c.DefaultQuery("direction", "asc")

	// Validate sort options
	validFields := map[string]bool{
		"id":       true,
		"type":     true,
		"category": true,
		"date":     true,
		"amount":   true,
	}

	if !validFields[field] {
		lg.Warn().Str("field", field).Msgf("%s Invalid sort field", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid sort field. Must be one of: id, type, category, date, amount",
		})
		return
	}

	if direction != "asc" && direction != "desc" {
		lg.Warn().Str("direction", direction).Msgf("%s Invalid sort direction", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid sort direction. Must be 'asc' or 'desc'",
		})
		return
	}

	sortOptions := model.SortOptions{
		Field:     field,
		Direction: direction,
	}

	records, err := h.sv.GetSalesRecords(c.Request.Context(), sortOptions)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to get sales records")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales records",
		})
		return
	}

	// Convert model.Data to SalesRecord DTO
	var responseRecords []SalesRecord
	for _, record := range records {
		responseRecords = append(responseRecords, SalesRecord{
			ID:       record.ID,
			Type:     record.Type,
			Category: record.Category,
			Date:     record.Date,
			Amount:   record.Amount,
		})
	}

	response := GetSalesRecordsResponse{
		Records: responseRecords,
	}

	lg.Debug().Int("count", len(responseRecords)).Msgf("%s records retrieved successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, response)
}

// UpdateSalesRecord handles PUT /items/:id for updating a record
func (h *SalesHandler) UpdateSalesRecord(c *ginext.Context) {
	lg := h.lg.With().Str("method", "UpdateSalesRecord").Logger()

	idParam := c.Param("id")
	id, err := parseID(idParam)
	if err != nil {
		lg.Warn().Str("id", idParam).Msgf("%s Invalid ID format", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	var request UpdateSalesRecordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		lg.Warn().Err(err).Msgf("%s failed to bind JSON", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	// Validate request fields
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

	err = h.sv.UpdateSalesRecord(c.Request.Context(), id, serviceRequest)
	if err != nil {
		lg.Error().Err(err).Int("ID", id).Msg("Failed to update sales record")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update sales record",
		})
		return
	}

	response := UpdateSalesRecordResponse{
		ID: id,
	}

	lg.Debug().Int("ID", id).Msgf("%s record updated successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, response)
}

// DeleteSalesRecord handles DELETE /items/:id for deleting a record
func (h *SalesHandler) DeleteSalesRecord(c *ginext.Context) {
	lg := h.lg.With().Str("method", "DeleteSalesRecord").Logger()

	idParam := c.Param("id")
	id, err := parseID(idParam)
	if err != nil {
		lg.Warn().Str("id", idParam).Msgf("%s Invalid ID format", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	err = h.sv.DeleteSalesRecord(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Int("ID", id).Msg("Failed to delete sales record")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete sales record",
		})
		return
	}

	response := DeleteSalesRecordResponse{
		ID: id,
	}

	lg.Debug().Int("ID", id).Msgf("%s record deleted successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, response)
}

// parseID is a helper function to parse ID from string to int
func parseID(idParam string) (int, error) {
	if idParam == "" {
		return 0, fmt.Errorf("empty ID")
	}

	id := 0
	for _, char := range idParam {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid ID format")
		}
		digit := int(char - '0')
		id = id*10 + digit
	}

	return id, nil
}

// GetAnalytics handles GET /analytics for retrieving analytics data
func (h *SalesHandler) GetAnalytics(c *ginext.Context) {
	lg := h.lg.With().Str("method", "GetAnalytics").Logger()

	// Get date range parameters from query
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	// Validate date parameters
	if from == "" {
		lg.Warn().Msgf("%s 'from' date parameter is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Parameter 'from' is required",
		})
		return
	}

	if to == "" {
		lg.Warn().Msgf("%s 'to' date parameter is required", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Parameter 'to' is required",
		})
		return
	}

	// Call service to get analytics data
	analytics, err := h.sv.GetAnalytics(c.Request.Context(), from, to)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to get analytics data")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get analytics data",
		})
		return
	}

	response := AnalyticsData{
		Sum:          analytics.Sum,
		Avg:          analytics.Avg,
		Count:        analytics.Count,
		Median:       analytics.Median,
		Percentile90: analytics.Percentile90,
	}

	lg.Debug().
		Float64("sum", analytics.Sum).
		Float64("avg", analytics.Avg).
		Int("count", analytics.Count).
		Float64("median", analytics.Median).
		Float64("percentile90", analytics.Percentile90).
		Msgf("%s analytics retrieved successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, response)
}

// ExportCSV handles GET /items/export for exporting sales records to CSV
func (h *SalesHandler) ExportCSV(c *ginext.Context) {
	lg := h.lg.With().Str("method", "ExportCSV").Logger()

	// Get date range parameters from query
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	// Call service to get CSV data
	csvData, err := h.sv.ExportCSV(c.Request.Context(), from, to)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to export CSV")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to export CSV",
		})
		return
	}

	// Set headers for CSV file download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\"sales_export.csv\"")

	lg.Debug().Int("size", len(csvData)).Msgf("%s CSV exported successfully", pkgConst.OpSuccess)

	c.Data(http.StatusOK, "text/csv", csvData)
}
