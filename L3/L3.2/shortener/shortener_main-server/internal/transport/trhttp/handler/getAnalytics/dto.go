package getAnalytics

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"

type response struct {
	TotalClicks int              `json:"total_clicks,omitempty"`
	Clicks      []model.Analitic `json:"clicks,omitempty"`
	Error       string           `json:"error,omitempty"`
}
