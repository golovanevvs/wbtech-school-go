package mainHandlers

type uploadResponse struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type imageResponse struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	OriginalPath  string `json:"original_path,omitempty"`
	ProcessedPath string `json:"processed_path,omitempty"`
	CreatedAt     string `json:"created_at"`
	Error         string `json:"error,omitempty"`
}

type deleteResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
