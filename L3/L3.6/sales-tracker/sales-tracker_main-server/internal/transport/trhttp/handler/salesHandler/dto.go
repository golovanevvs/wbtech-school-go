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
