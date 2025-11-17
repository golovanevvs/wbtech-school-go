package model

import "time"

type ImageStatus string

const (
	StatusUploading  ImageStatus = "uploading"
	StatusProcessing ImageStatus = "processing"
	StatusCompleted  ImageStatus = "completed"
	StatusFailed     ImageStatus = "failed"
)

type Image struct {
	ID            int
	Status        ImageStatus
	OriginalPath  string
	ProcessedPath *string
	CreatedAt     time.Time
}
