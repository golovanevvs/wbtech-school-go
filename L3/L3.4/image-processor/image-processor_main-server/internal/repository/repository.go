package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/repository/rpFileStorage"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
	*rpFileStorage.RpFileStorage
}

func New(cfg *Config, pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) (*Repository, error) {
	return &Repository{
		RpPostgres:    rpPostgres.New(pg, rs),
		RpFileStorage: rpFileStorage.New(cfg.RpFileStorage),
	}, nil

}
