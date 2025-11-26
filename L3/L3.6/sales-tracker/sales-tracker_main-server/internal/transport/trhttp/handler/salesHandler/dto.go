package salesHandler

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateSalesRecordRequest represents request for creating sales record
type CreateSalesRecordRequest struct {
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Amount   float64 `json:"amount"`
}

// CreateSalesRecordResponse represents response after creating sales record
type CreateSalesRecordResponse struct {
	ID int `json:"id"`
}

// GetSalesRecordsRequest represents request for getting sales records with sorting
type GetSalesRecordsRequest struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// SalesRecord represents a sales record in responses
type SalesRecord struct {
	ID       int     `json:"id"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Amount   float64 `json:"amount"`
}

// GetSalesRecordsResponse represents response for getting sales records
type GetSalesRecordsResponse struct {
	Records []SalesRecord `json:"records"`
}

// UpdateSalesRecordRequest represents request for updating sales record
type UpdateSalesRecordRequest struct {
	Type     string  `json:"type"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Amount   float64 `json:"amount"`
}

// UpdateSalesRecordResponse represents response after updating sales record
type UpdateSalesRecordResponse struct {
	ID int `json:"id"`
}

// DeleteSalesRecordResponse represents response after deleting sales record
type DeleteSalesRecordResponse struct {
	ID int `json:"id"`
}

// AnalyticsRequest represents request for getting analytics
type AnalyticsRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// AnalyticsData represents analytics metrics response
type AnalyticsData struct {
	Sum          float64 `json:"sum"`
	Avg          float64 `json:"avg"`
	Count        int     `json:"count"`
	Median       float64 `json:"median"`
	Percentile90 float64 `json:"percentile90"`
}
