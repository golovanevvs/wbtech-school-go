package mainHandlers

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"

type uploadResponse struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type imageResponse struct {
	ID           int                  `json:"id"`
	Status       string               `json:"status"`
	OriginalPath string               `json:"original_path,omitempty"`
	ProcessedUrl string               `json:"processed_url,omitempty"`
	CreatedAt    string               `json:"created_at"`
	Operations   model.ProcessOptions `json:"operations,omitempty"`
	Error        string               `json:"error,omitempty"`
}

type deleteResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
