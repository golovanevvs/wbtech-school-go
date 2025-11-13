package repository_test

import (
	"testing"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		redisClient *pkgRedis.Client
		expectErr   bool
	}{
		{
			name:        "valid redis client",
			redisClient: &pkgRedis.Client{},
			expectErr:   false,
		},
		{
			name:        "nil redis client",
			redisClient: nil,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := repository.New(tt.redisClient)

			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, repo)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, repo)
			require.NotNil(t, repo.RpRedis)
		})
	}
}
