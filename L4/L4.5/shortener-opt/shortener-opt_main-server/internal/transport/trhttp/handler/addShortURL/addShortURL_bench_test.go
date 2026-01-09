package addShortURL

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type MockService struct{}

func (m *MockService) AddShortURL(ctx context.Context, originalURL, short string) (int, string, error) {
	return 1, "abc123", nil
}

func setupTestHandler() (*Handler, *zlog.Zerolog) {
	zlog.InitConsole()
	lg := zlog.Logger.With().Str("component", "test").Logger()
	rt := ginext.New("test")
	h := New(&lg, rt, &MockService{})
	h.RegisterRoutes()
	return h, &lg
}

func BenchmarkAddShortURL(b *testing.B) {
	h, _ := setupTestHandler()

	body := map[string]string{
		"url": "https://example.com/very/long/url/that/needs/to/be/shortened/and/test/performance",
	}
	jsonBody, _ := json.Marshal(body)
	reader := bytes.NewReader(jsonBody)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		req := httptest.NewRequest("POST", "/shorten", reader)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.rt.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Errorf("unexpected status: %d", w.Code)
		}
	}
}

func BenchmarkGetOriginalURL(b *testing.B) {
	h, _ := setupTestHandler()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/abc123", nil)
		w := httptest.NewRecorder()
		h.rt.ServeHTTP(w, req)
	}
}
