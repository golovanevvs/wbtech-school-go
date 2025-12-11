package addShortURL

type request struct {
	Original string `json:"original"`
	Short    string `json:"short"`
}

type response struct {
	Short string `json:"short,omitempty"`
	Error string `json:"error,omitempty"`
}
