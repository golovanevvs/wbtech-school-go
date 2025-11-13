package addShortURL

type request struct {
	Original string `json:"original"`
	Short    string `json:"short"`
}

type response struct {
	ShortURL string `json:"short_url,omitempty"`
	Error    string `json:"error,omitempty"`
}
