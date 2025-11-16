package mainHandlers

type request struct {
	ParentID *int   `json:"parent_id"`
	Text     string `json:"text" binding:"required"`
}

type response struct {
	ID        int         `json:"id,omitempty"`
	ParentID  *int        `json:"parent_id,omitempty"`
	Text      string      `json:"text,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt *string     `json:"updated_at,omitempty"`
	Children  []*response `json:"children,omitempty"`
	Error     string      `json:"error,omitempty"`
}
