package addShortURL

import (
	"context"
	"testing"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/model"
)

type mockRp struct {
	saveCount int
}

func (m *mockRp) SaveShortURL(ctx context.Context, shortURL model.ShortURL) (id int, err error) {
	m.saveCount++
	return m.saveCount, nil
}

func BenchmarkAddShortURLService(b *testing.B) {
	repo := &mockRp{}
	svc := New(repo)

	ctx := context.Background()
	testURL := "https://example.com/very/long/url/that/needs/to/be/shortened/and/test/performance/optimization"

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, err := svc.AddShortURL(ctx, testURL, "")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateShortCode(b *testing.B) {
	svc := New(&mockRp{})

	// Subtest for old version (uuid)
	b.Run("old", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			svc.generateShortCodeOld()
		}
	})

	// Subtest for new version (sync.Pool)
	b.Run("new", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			svc.generateShortCode()
		}
	})
}

func BenchmarkAddShortURLServiceWithCustomShort(b *testing.B) {
	repo := &mockRp{}
	svc := New(repo)

	ctx := context.Background()
	testURL := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, err := svc.AddShortURL(ctx, testURL, "custom123")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestAddShortURLTrace(t *testing.T) {
	repo := &mockRp{}
	svc := New(repo)

	ctx := context.Background()
	testURL := "https://example.com/test"

	for range 1000 {
		_, _, err := svc.AddShortURL(ctx, testURL, "")
		if err != nil {
			t.Fatal(err)
		}
	}
}
