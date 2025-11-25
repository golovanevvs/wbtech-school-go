package model

// Data represents a sales/expense record
type Data struct {
	ID       int
	Type     string
	Category string
	Date     string
	Amount   float64
}

// Analytics represents analytics metrics
type Analytics struct {
	Sum          float64
	Avg          float64
	Count        int
	Median       float64
	Percentile90 float64
}

// SortOptions represents sorting parameters
type SortOptions struct {
	Field     string
	Direction string
}
